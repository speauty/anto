import 'dart:io';

/// Small Windows-only bridge for the native file and folder dialogs.
///
/// The desktop client is currently Windows-only. Keeping this bridge behind a
/// class makes the UI independent from the dialog implementation and leaves a
/// straightforward seam for widget tests and a future cross-platform picker.
class WindowsPicker {
  const WindowsPicker();

  Future<List<String>> pickFiles() async {
    if (!Platform.isWindows) return const <String>[];
    const script = r'''
Add-Type -AssemblyName System.Windows.Forms
[System.Windows.Forms.Application]::EnableVisualStyles()
$dialog = New-Object System.Windows.Forms.OpenFileDialog
$dialog.Multiselect = $true
$dialog.CheckFileExists = $true
$dialog.Filter = 'SubRip subtitle (*.srt)|*.srt'
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
  $dialog.FileNames | ForEach-Object { Write-Output $_ }
}
''';
    final result = await Process.run(
        'powershell.exe', <String>['-NoProfile', '-STA', '-Command', script]);
    if (result.exitCode != 0) {
      throw StateError('打开文件选择器失败: ${result.stderr}'.trim());
    }
    return _lines(result.stdout.toString());
  }

  Future<String?> pickDirectory() async {
    if (!Platform.isWindows) return null;
    const script = r'''
Add-Type -AssemblyName System.Windows.Forms
[System.Windows.Forms.Application]::EnableVisualStyles()
$dialog = New-Object System.Windows.Forms.FolderBrowserDialog
$dialog.Description = '选择包含 SRT 文件的目录'
$dialog.ShowNewFolderButton = $false
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
  Write-Output $dialog.SelectedPath
}
''';
    final result = await Process.run(
        'powershell.exe', <String>['-NoProfile', '-STA', '-Command', script]);
    if (result.exitCode != 0) {
      throw StateError('打开目录选择器失败: ${result.stderr}'.trim());
    }
    final lines = _lines(result.stdout.toString());
    return lines.isEmpty ? null : lines.first;
  }

  List<String> _lines(String output) => output
      .split(RegExp(r'\r?\n'))
      .map((line) => line.trim())
      .where((line) => line.isNotEmpty)
      .toList();
}
