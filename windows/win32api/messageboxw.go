//go:build windows

package messageboxw

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func WarningManyPortsNotAvailable(portMin int, portMax int) {
	// https://stackoverflow.com/questions/46705163/how-to-alert-in-go-to-show-messagebox
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
	var user32DLL = syscall.NewLazyDLL("user32.dll")
	var procMessageBox = user32DLL.NewProc("MessageBoxW") // Return value: Type int

	const (
		MB_OK           = 0x00000000
		MB_OKCANCEL     = 0x00000001
		MB_YESNO        = 0x00000004
		MB_SYSTEMMODAL  = 0x00001000
		MB_ICONQUESTION = 0x00000020
		MB_ICONWARNING  = 0x00000030
		MB_ICONASTERISK = 0x00000040
	)

	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw#return-value
	lpCaption, _ := syscall.UTF16PtrFromString("Opening Application might take some time")                                                                                                                                                                                                                                            // LPCWSTR
	lpText, _ := syscall.UTF16PtrFromString("Many ports between " + strconv.Itoa(portMin) + " - " + strconv.Itoa(portMax) + " already used\n\nDo you want to continue searching for ports to open program? \n\nThis might take 30 seconds from pressing Yes. \nPressing No will stop application from launching.\n\nPress Yes or No") // LPCWSTR

	buttonPressed, _, _ := syscall.SyscallN(procMessageBox.Addr(),
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		MB_YESNO|MB_ICONQUESTION, // Let the window TOPMOST.
	)

	const yes = 6
	const no = 7

	if buttonPressed == yes {
		fmt.Println("Continuing search for open ports")
	} else if buttonPressed == no {
		fmt.Println("Closing application because of cancel")
		os.Exit(0)
	}

}

func WarningNoPortsAvailable() {
	// https://stackoverflow.com/questions/46705163/how-to-alert-in-go-to-show-messagebox
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
	var user32DLL = syscall.NewLazyDLL("user32.dll")
	var procMessageBox = user32DLL.NewProc("MessageBoxW") // Return value: Type int
	const (
		MB_OK           = 0x00000000
		MB_OKCANCEL     = 0x00000001
		MB_YESNO        = 0x00000004
		MB_SYSTEMMODAL  = 0x00001000
		MB_ICONQUESTION = 0x00000020
		MB_ICONWARNING  = 0x00000030
		MB_ICONASTERISK = 0x00000040
	)

	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw#return-value
	lpCaption, _ := syscall.UTF16PtrFromString("No ports")                                                                // LPCWSTR
	lpText, _ := syscall.UTF16PtrFromString("There was no available ports to launch the application, stopping launch...") // LPCWSTR

	_, _, _ = syscall.SyscallN(procMessageBox.Addr(),
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		MB_OK|MB_ICONWARNING, // Let the window TOPMOST.
	)
}

func WarningYouNeedToInstallChrome() {
	// https://stackoverflow.com/questions/46705163/how-to-alert-in-go-to-show-messagebox
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
	var user32DLL = syscall.NewLazyDLL("user32.dll")
	var procMessageBox = user32DLL.NewProc("MessageBoxW") // Return value: Type int
	const (
		MB_OK           = 0x00000000
		MB_OKCANCEL     = 0x00000001
		MB_YESNO        = 0x00000004
		MB_SYSTEMMODAL  = 0x00001000
		MB_ICONQUESTION = 0x00000020
		MB_ICONWARNING  = 0x00000030
		MB_ICONASTERISK = 0x00000040
	)

	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw#return-value
	lpCaption, _ := syscall.UTF16PtrFromString("This program has a dependency")                                                                                                                                                                                                                                      // LPCWSTR
	lpText, _ := syscall.UTF16PtrFromString("You need to install Google Chrome in order to use the program KeyFiles developed by Beksoft ApS\n\nYou can download Chrome at: https://www.google.com/chrome/\n\nPlease install Chrome at: " + os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe") // LPCWSTR

	_, _, _ = syscall.SyscallN(procMessageBox.Addr(),
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		MB_OK|MB_ICONWARNING, // Let the window TOPMOST.
	)
}
