#include "flutter_window.h"

#include <shellapi.h>

#include <optional>

#include <string>
#include <cstdint>
#include <limits>
#include <vector>

#include "resource.h"

namespace {
constexpr UINT kTrayCallbackMessage = WM_APP + 1;
constexpr UINT kTrayIconId = 1;
constexpr UINT kTrayExitCommand = 1;
}  // namespace

#include "flutter/generated_plugin_registrant.h"

FlutterWindow::FlutterWindow(const flutter::DartProject& project)
    : project_(project) {}

FlutterWindow::~FlutterWindow() { CloseCoreJob(); }

bool FlutterWindow::AttachCoreProcess(DWORD process_id) {
  if (!core_job_) {
    core_job_ = CreateJobObjectW(nullptr, nullptr);
    if (!core_job_) return false;
    JOBOBJECT_EXTENDED_LIMIT_INFORMATION limits{};
    limits.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE;
    if (!SetInformationJobObject(core_job_, JobObjectExtendedLimitInformation,
                                 &limits, sizeof(limits))) {
      CloseCoreJob();
      return false;
    }
  }
  HANDLE process = OpenProcess(PROCESS_SET_QUOTA | PROCESS_TERMINATE, FALSE,
                               process_id);
  if (!process) return false;
  const bool attached = AssignProcessToJobObject(core_job_, process) == TRUE;
  CloseHandle(process);
  return attached;
}

void FlutterWindow::CloseCoreJob() {
  if (!core_job_) return;
  CloseHandle(core_job_);
  core_job_ = nullptr;
}

bool FlutterWindow::OnCreate() {
  if (!Win32Window::OnCreate()) {
    return false;
  }

  RECT frame = GetClientArea();

  // The size here must match the window dimensions to avoid unnecessary surface
  // creation / destruction in the startup path.
  flutter_controller_ = std::make_unique<flutter::FlutterViewController>(
      frame.right - frame.left, frame.bottom - frame.top, project_);
  // Ensure that basic setup of the controller was successful.
  if (!flutter_controller_->engine() || !flutter_controller_->view()) {
    return false;
  }
  RegisterPlugins(flutter_controller_->engine());
  SetChildContent(flutter_controller_->view()->GetNativeWindow());
  DragAcceptFiles(GetHandle(), TRUE);
  drop_channel_ = std::make_unique<flutter::MethodChannel<flutter::EncodableValue>>(
      flutter_controller_->engine()->messenger(), "anto/windows_drop",
      &flutter::StandardMethodCodec::GetInstance());
  window_channel_ = std::make_unique<
      flutter::MethodChannel<flutter::EncodableValue>>(
      flutter_controller_->engine()->messenger(), "anto/windows_window",
      &flutter::StandardMethodCodec::GetInstance());
  window_channel_->SetMethodCallHandler(
      [this](const auto& call, auto result) {
        if (call.method_name() == "getAppVersion") {
          result->Success(flutter::EncodableValue(std::string(FLUTTER_VERSION)));
          return;
        }
        if (call.method_name() == "setMinimizeToTray") {
          if (call.arguments() &&
              std::holds_alternative<bool>(*call.arguments())) {
            minimize_to_tray_ = std::get<bool>(*call.arguments());
          }
          result->Success();
          return;
        }
        if (call.method_name() == "attachCoreProcess") {
          int64_t process_id = 0;
          if (call.arguments() &&
              std::holds_alternative<int32_t>(*call.arguments())) {
            process_id = std::get<int32_t>(*call.arguments());
          } else if (call.arguments() &&
                     std::holds_alternative<int64_t>(*call.arguments())) {
            process_id = std::get<int64_t>(*call.arguments());
          }
          if (process_id <= 0 ||
              process_id > std::numeric_limits<DWORD>::max()) {
            result->Error("invalid_process_id", "Go Core process ID is invalid");
            return;
          }
          if (!AttachCoreProcess(static_cast<DWORD>(process_id))) {
            result->Error("core_job_failed",
                          "Failed to attach Go Core to the UI process");
            return;
          }
          result->Success();
          return;
        }
        result->NotImplemented();
      });

  flutter_controller_->engine()->SetNextFrameCallback([&]() {
    this->Show();
  });

  // Flutter can complete the first frame before the "show window" callback is
  // registered. The following call ensures a frame is pending to ensure the
  // window is shown. It is a no-op if the first frame hasn't completed yet.
  flutter_controller_->ForceRedraw();

  return true;
}

void FlutterWindow::OnDestroy() {
  CloseCoreJob();
  if (flutter_controller_) {
    flutter_controller_ = nullptr;
  }

  Win32Window::OnDestroy();
}

LRESULT
FlutterWindow::MessageHandler(HWND hwnd, UINT const message,
                              WPARAM const wparam,
                              LPARAM const lparam) noexcept {
  // Give Flutter, including plugins, an opportunity to handle window messages.
  if (flutter_controller_) {
    std::optional<LRESULT> result =
        flutter_controller_->HandleTopLevelWindowProc(hwnd, message, wparam,
                                                      lparam);
    if (result) {
      return *result;
    }
  }

  switch (message) {
    case WM_CLOSE:
      if (minimize_to_tray_ && !exiting_) {
        ShowTrayIcon();
        ShowWindow(hwnd, SW_HIDE);
        return 0;
      }
      break;
    case kTrayCallbackMessage:
      if (wparam != kTrayIconId) break;
      if (lparam == WM_LBUTTONUP || lparam == WM_LBUTTONDBLCLK) {
        RestoreFromTray();
        return 0;
      }
      if (lparam == WM_RBUTTONUP) {
        POINT point;
        GetCursorPos(&point);
        HMENU menu = CreatePopupMenu();
        if (menu) {
          AppendMenuW(menu, MF_STRING, kTrayExitCommand, L"Exit");
          SetForegroundWindow(hwnd);
          const auto command = TrackPopupMenu(
              menu, TPM_RETURNCMD | TPM_NONOTIFY, point.x, point.y, 0, hwnd,
              nullptr);
          DestroyMenu(menu);
          if (command == kTrayExitCommand) ExitFromTray();
        }
        return 0;
      }
      break;
    case WM_DROPFILES: {
      const auto drop = reinterpret_cast<HDROP>(wparam);
      const auto count = DragQueryFileW(drop, 0xFFFFFFFF, nullptr, 0);
      flutter::EncodableList paths;
      for (UINT index = 0; index < count; ++index) {
        const auto length = DragQueryFileW(drop, index, nullptr, 0);
        std::vector<wchar_t> buffer(length + 1);
        DragQueryFileW(drop, index, buffer.data(), length + 1);
        const auto utf8_length = WideCharToMultiByte(
            CP_UTF8, 0, buffer.data(), static_cast<int>(length), nullptr, 0,
            nullptr, nullptr);
        std::string path(utf8_length, '\0');
        WideCharToMultiByte(CP_UTF8, 0, buffer.data(), static_cast<int>(length),
                            path.data(), utf8_length, nullptr, nullptr);
        paths.emplace_back(path);
      }
      DragFinish(drop);
      if (drop_channel_ && !paths.empty()) {
        drop_channel_->InvokeMethod(
            "filesDropped",
            std::make_unique<flutter::EncodableValue>(paths));
      }
      return 0;
    }
    case WM_FONTCHANGE:
      flutter_controller_->engine()->ReloadSystemFonts();
      break;
  }

  return Win32Window::MessageHandler(hwnd, message, wparam, lparam);
}

void FlutterWindow::ShowTrayIcon() {
  if (tray_visible_) return;
  NOTIFYICONDATAW data{};
  data.cbSize = sizeof(data);
  data.hWnd = GetHandle();
  data.uID = kTrayIconId;
  data.uFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP;
  data.uCallbackMessage = kTrayCallbackMessage;
  data.hIcon = LoadIcon(GetModuleHandle(nullptr),
                        MAKEINTRESOURCE(IDI_APP_ICON));
  wcscpy_s(data.szTip, L"anto");
  tray_visible_ = Shell_NotifyIconW(NIM_ADD, &data) == TRUE;
}

void FlutterWindow::HideTrayIcon() {
  if (!tray_visible_) return;
  NOTIFYICONDATAW data{};
  data.cbSize = sizeof(data);
  data.hWnd = GetHandle();
  data.uID = kTrayIconId;
  Shell_NotifyIconW(NIM_DELETE, &data);
  tray_visible_ = false;
}

void FlutterWindow::RestoreFromTray() {
  HideTrayIcon();
  ShowWindow(GetHandle(), SW_SHOW);
  SetForegroundWindow(GetHandle());
}

void FlutterWindow::ExitFromTray() {
  exiting_ = true;
  HideTrayIcon();
  DestroyWindow(GetHandle());
}
