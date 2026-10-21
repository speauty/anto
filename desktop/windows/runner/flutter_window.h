#ifndef RUNNER_FLUTTER_WINDOW_H_
#define RUNNER_FLUTTER_WINDOW_H_

#include <flutter/dart_project.h>
#include <flutter/method_channel.h>
#include <flutter/standard_method_codec.h>
#include <flutter/flutter_view_controller.h>

#include <memory>

#include "win32_window.h"

// A window that does nothing but host a Flutter view.
class FlutterWindow : public Win32Window {
 public:
  // Creates a new FlutterWindow hosting a Flutter view running |project|.
  explicit FlutterWindow(const flutter::DartProject& project);
  virtual ~FlutterWindow();

 protected:
  // Win32Window:
  bool OnCreate() override;
  void OnDestroy() override;
  LRESULT MessageHandler(HWND window, UINT const message, WPARAM const wparam,
                         LPARAM const lparam) noexcept override;

 private:
  // The project to run.
  flutter::DartProject project_;

  // The Flutter instance hosted by this window.
  std::unique_ptr<flutter::FlutterViewController> flutter_controller_;

  // Delivers native Windows file drops to the Flutter client.
  std::unique_ptr<flutter::MethodChannel<flutter::EncodableValue>> drop_channel_;

  // Controls the optional Windows notification-area behavior.
  std::unique_ptr<flutter::MethodChannel<flutter::EncodableValue>>
      window_channel_;
  bool minimize_to_tray_ = false;
  bool tray_visible_ = false;
  bool exiting_ = false;
  HANDLE core_job_ = nullptr;

  bool AttachCoreProcess(DWORD process_id);
  void CloseCoreJob();

  void ShowTrayIcon();
  void HideTrayIcon();
  void RestoreFromTray();
  void ExitFromTray();
};

#endif  // RUNNER_FLUTTER_WINDOW_H_
