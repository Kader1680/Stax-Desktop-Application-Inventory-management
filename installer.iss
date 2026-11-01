
#define MyAppName "STAX"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "STAX"
#define MyAppExeName "stax-desktop.exe"

[Setup]
AppId={{A82D91F4-63C7-4B0E-9F11-123456789ABC}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\STAX
DefaultGroupName=STAX
OutputDir=installer-output
OutputBaseFilename=STAX-Setup-1.0.0
SetupIconFile=build\windows\icon.ico
ArchitecturesInstallIn64BitMode=x64
PrivilegesRequired=admin
Compression=lzma
SolidCompression=yes
WizardStyle=modern

[Tasks]
Name: "desktopicon"; Description: "Create a desktop shortcut"; Flags: unchecked

[Files]
Source: "build\bin\stax-desktop.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\STAX"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\STAX"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "Launch STAX"; Flags: postinstall nowait skipifsilent
