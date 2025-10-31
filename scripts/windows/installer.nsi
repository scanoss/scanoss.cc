; SCANOSS Code Compare NSIS Installer Script
; This script creates a Windows installer with optional PATH management

!include "MUI2.nsh"
!include "x64.nsh"
!include "WinMessages.nsh"

; Application Settings
!define APP_NAME "SCANOSS Code Compare"
!define APP_EXE "scanoss-cc-windows.exe"
!define APP_CLI_NAME "scanoss-cc"
!define COMPANY_NAME "SCANOSS"
!define INSTALL_DIR_NAME "SCANOSS"
!define REG_UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"

; Installer Settings
Name "${APP_NAME}"
OutFile "..\..\dist\${APP_NAME}-Setup.exe"
InstallDir "$PROGRAMFILES64\${INSTALL_DIR_NAME}"
InstallDirRegKey HKLM "Software\${INSTALL_DIR_NAME}" "Install_Dir"
RequestExecutionLevel admin
SilentInstall normal

; Modern UI Configuration
!define MUI_ABORTWARNING
;!define MUI_ICON "..\..\assets\appicon.ico"  ; Optional: add if you have an .ico file
;!define MUI_UNICON "..\..\assets\appicon.ico"  ; Optional: add if you have an .ico file
;!define MUI_WELCOMEFINISHPAGE_BITMAP "..\..\assets\installer_banner.bmp"  ; Optional: add if you have a banner

; Pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\..\LICENSE"
!insertmacro MUI_PAGE_DIRECTORY

; Custom page for PATH selection
Page custom PathOptionsPage

!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Languages
!insertmacro MUI_LANGUAGE "English"

; Variables
Var AddToPath
Var AddToPathCheckbox

; Helper function to add to PATH
!define Environ 'HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"'

Function AddToPath
  Exch $0
  Push $1
  Push $2
  Push $3

  ; Read current PATH
  ReadRegStr $1 ${Environ} "PATH"

  ; Check if already in PATH
  Push "$1;"
  Push "$0;"
  Call StrStr
  Pop $2
  StrCmp $2 "" 0 done

  ; Append to PATH
  StrCpy $2 "$1;$0"
  WriteRegExpandStr ${Environ} "PATH" $2

  ; Broadcast environment change
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

  done:
  Pop $3
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Function un.RemoveFromPath
  Exch $0
  Push $1
  Push $2
  Push $3
  Push $4
  Push $5
  Push $6

  ; Read current PATH
  ReadRegStr $1 ${Environ} "PATH"

  StrCpy $5 $1 1 -1 ; Check if last char is ;
  ${If} $5 != ";"
    StrCpy $1 "$1;" ; Add ; to end if not present
  ${EndIf}

  Push $1
  Push "$0;"
  Call un.StrStr
  Pop $2

  StrCmp $2 "" done
  StrLen $3 "$0;"
  StrLen $4 $2
  StrCpy $5 $1 -$4
  StrCpy $6 $2 "" $3
  StrCpy $3 "$5$6"

  ; Remove trailing ;
  StrCpy $5 $3 1 -1
  ${If} $5 == ";"
    StrCpy $3 $3 -1
  ${EndIf}

  WriteRegExpandStr ${Environ} "PATH" $3

  ; Broadcast environment change
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

  done:
  Pop $6
  Pop $5
  Pop $4
  Pop $3
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

; String search function
Function StrStr
  Exch $R1 ; needle
  Exch
  Exch $R2 ; haystack
  Push $R3
  Push $R4
  Push $R5

  StrLen $R3 $R1
  StrCpy $R4 0

  loop:
    StrCpy $R5 $R2 $R3 $R4
    StrCmp $R5 $R1 done
    StrCmp $R5 "" done
    IntOp $R4 $R4 + 1
    Goto loop

  done:
    StrCpy $R1 $R2 "" $R4
    Pop $R5
    Pop $R4
    Pop $R3
    Pop $R2
    Exch $R1
FunctionEnd

Function un.StrStr
  Exch $R1
  Exch
  Exch $R2
  Push $R3
  Push $R4
  Push $R5

  StrLen $R3 $R1
  StrCpy $R4 0

  loop:
    StrCpy $R5 $R2 $R3 $R4
    StrCmp $R5 $R1 done
    StrCmp $R5 "" done
    IntOp $R4 $R4 + 1
    Goto loop

  done:
    StrCpy $R1 $R2 "" $R4
    Pop $R5
    Pop $R4
    Pop $R3
    Pop $R2
    Exch $R1
FunctionEnd

; Custom PATH Options Page
Function PathOptionsPage
  !insertmacro MUI_HEADER_TEXT "Installation Options" "Choose additional options"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateLabel} 0 0 100% 20u "Select additional installation options:"
  Pop $0

  ${NSD_CreateCheckBox} 10 30u 100% 10u "Add scanoss-cc to system PATH (recommended)"
  Pop $AddToPathCheckbox
  ${NSD_Check} $AddToPathCheckbox  ; Checked by default

  ${NSD_CreateLabel} 20 45u 100% 20u "This allows you to run 'scanoss-cc' from any command prompt"
  Pop $0

  nsDialogs::Show
FunctionEnd

; Store PATH choice
Function PathOptionsPageLeave
  ${NSD_GetState} $AddToPathCheckbox $AddToPath
FunctionEnd

; Main Install Section
Section "Install"
  SetOutPath "$INSTDIR"

  ; Copy application files
  File /r "..\..\build\bin\${APP_EXE}"
  File /r "..\..\build\assets"
  File "..\..\README.md"
  File "..\..\LICENSE"

  ; Create start menu shortcuts
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe"

  ; Desktop shortcut
  CreateShortcut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}"

  ; Add to PATH if selected
  ${If} $AddToPath == ${BST_CHECKED}
    ; Add to system PATH for all users
    Push "$INSTDIR"
    Call AddToPath
  ${EndIf}

  ; Register App Paths for command invocation
  WriteRegStr HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\${APP_CLI_NAME}.exe" "" "$INSTDIR\${APP_EXE}"
  WriteRegStr HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\${APP_CLI_NAME}.exe" "Path" "$INSTDIR"

  ; Write installation info
  WriteRegStr HKLM "Software\${INSTALL_DIR_NAME}" "Install_Dir" "$INSTDIR"
  WriteRegStr HKLM "Software\${INSTALL_DIR_NAME}" "Version" "${APP_VERSION}"

  ; Write uninstall information
  WriteRegStr HKLM "${REG_UNINSTALL_KEY}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "${REG_UNINSTALL_KEY}" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr HKLM "${REG_UNINSTALL_KEY}" "DisplayIcon" "$INSTDIR\${APP_EXE}"
  WriteRegStr HKLM "${REG_UNINSTALL_KEY}" "Publisher" "${COMPANY_NAME}"
  WriteRegStr HKLM "${REG_UNINSTALL_KEY}" "DisplayVersion" "${APP_VERSION}"
  WriteRegDWORD HKLM "${REG_UNINSTALL_KEY}" "NoModify" 1
  WriteRegDWORD HKLM "${REG_UNINSTALL_KEY}" "NoRepair" 1

  ; Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"

SectionEnd

; Uninstall Section
Section "Uninstall"
  ; Remove from PATH if it was added
  Push "$INSTDIR"
  Call un.RemoveFromPath

  ; Remove App Paths registry
  DeleteRegKey HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\${APP_CLI_NAME}.exe"

  ; Remove files and directories
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\Uninstall.exe"
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\LICENSE"
  RMDir /r "$INSTDIR\assets"
  RMDir "$INSTDIR"

  ; Remove shortcuts
  Delete "$SMPROGRAMS\${APP_NAME}\*.*"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  ; Remove registry keys
  DeleteRegKey HKLM "${REG_UNINSTALL_KEY}"
  DeleteRegKey HKLM "Software\${INSTALL_DIR_NAME}"

SectionEnd
