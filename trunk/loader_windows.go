// +build windows
/**
 * trunk - Steam Modding Utility Package
 * https://github.com/awstanley/GoSteam/trunk
 *
 * Copyright (C) 2016 A.W. 'Swixel' Stanley <code@swixel.net>
 *
 * This software is provided 'as-is', without any express or
 * implied warranty. In no event will the authors be held
 * liable for any damages arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute
 * it freely, subject to the following restrictions:
 *
 *   1. The origin of this software must not be misrepresented;
 *      you must not claim that you wrote the original software.
 *      If you use this software in a product, an acknowledgment
 *      in the product documentation would be appreciated but is
 *      not required.
 *   2. Altered source versions must be plainly marked as such,
 *      and must not be misrepresented as being the original
 *      software.
 *   3. This notice may not be removed or altered from
 *      any source distribution.
**/

package trunk

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var (
	// This is a pointer to the Steam Library
	pSteamLibrary = syscall.NewLazyDLL(SteamLibraryName)

	// Init
	pSteamApiInit = pSteamLibrary.NewProc("SteamAPI_Init")

	pSteamApps = pSteamLibrary.NewProc("SteamApps")

	// GetappInstallDir (why we're here at the moment)
	pSteamApiISteamAppsGetAppInstallDir = pSteamLibrary.NewProc("SteamAPI_ISteamApps_GetAppInstallDir")
)

var stringBufferSize = 32 * 1024

// steamApiInit initialises the Steam connection
func steamApiInit() bool {
	r1, _, _ := pSteamApiInit.Call()
	return r1 == 1
}

// GetAppInstallDir returns the application directory (or "")
func GetAppInstallDir() string {

	ptr, _, _ := pSteamApps.Call()

	if ptr == 0 {
		log.Println("Failed on call to SteamApps() [pid is %p]\n", ptr)
		return ""
	}

	// Get AppId
	var appId uint32
	fmt.Sscanf(os.Getenv("SteamAppId"), "%d", &appId)

	if appId == 0 {
		log.Println("Failed to get AppID from env\n")
		return ""
	}

	// Allocate buffer
	buf := make([]byte, stringBufferSize)
	r1, _, _ := pSteamApiISteamAppsGetAppInstallDir.Call(ptr,
		uintptr(appId),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(stringBufferSize),
	)

	// Create the string
	str := string(buf[0 : uint32(r1)-1])

	// Return it
	return str
}
