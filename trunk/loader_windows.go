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

	// SteamApps
	pSteamApps = pSteamLibrary.NewProc("SteamApps")

	// GetappInstallDir (why we're here at the moment)
	pSteamAppsGetAppInstallDir = pSteamLibrary.NewProc("SteamAPI_ISteamApps_GetAppInstallDir")

	// SteamFriends
	pSteamFriends = pSteamLibrary.NewProc("SteamFriends")

	// Gets the persona name
	pISteamGetPersonaName = pSteamLibrary.NewProc("SteamAPI_ISteamFriends_GetPersonaName")

	// SteamApps
	pSteamUser = pSteamLibrary.NewProc("SteamUser")

	// SteamUser
	pSteamUserGetSteamID = pSteamLibrary.NewProc("SteamAPI_ISteamUser_GetSteamID")
)

var stringBufferSize = 32 * 1024

// steamApiInit initialises the Steam connection
func steamApiInit() bool {
	r1, _, _ := pSteamApiInit.Call()
	return r1 == 1
}

// GetAppInstallDir returns the application directory (or "")
func GetAppInstallDir() string {

	// Gets the "SteamApps" instance.
	ptr, _, _ := pSteamApps.Call()

	// Fails if it's nil/null (0)
	if ptr == 0 {
		log.Println("Failed on call to SteamApps() [pid is %p]\n", ptr)
		return ""
	}

	// Gets the AppId (previously set)
	var appId uint32
	fmt.Sscanf(os.Getenv("SteamAppId"), "%d", &appId)

	// Fails if it nil/null (0)
	if appId == 0 {
		log.Println("Failed to get AppID from env\n")
		return ""
	}

	// Allocates a HUGE string buffer to hold the path,
	// to handle the really weird cases people have.
	buf := make([]byte, stringBufferSize)

	// Calls the Steamworks "GetAppInstallDir" with the given App.
	r1, _, _ := pSteamAppsGetAppInstallDir.Call(ptr,
		uintptr(appId),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(stringBufferSize),
	)

	// Turns the buffer into a string of the path
	str := string(buf[0 : uint32(r1)-1])

	// Returns the path
	return str
}

// GetSteamID64 returns the 64-bit SteamID of the current user.
// This can be used to handle a few different cases (e.g. userdata)
func GetSteamID64() uint64 {

	// Gets the SteamUser
	ptr, _, _ := pSteamUser.Call()

	// Fail if it's nil/null/0
	if ptr == 0 {
		log.Println("Failed on call to SteamUser() [pid is %p]\n", ptr)
		return 0
	}

	// Get AppId from the environment
	var appId uint32
	fmt.Sscanf(os.Getenv("SteamAppId"), "%d", &appId)

	// Fail if it's nil/null/0
	if appId == 0 {
		log.Println("Failed to get AppID from env\n")
		return 0
	}

	// Make the call (this one's easy)
	r1, _, _ := pSteamUserGetSteamID.Call(ptr)

	// Return it
	return uint64(r1)
}

// GetPersonaName returns the current "friends" name
func GetPersonaName() string {

	// Gets the "SteamFriends" instance.
	ptr, _, _ := pSteamFriends.Call()

	// Fails if it's nil/null (0)
	if ptr == 0 {
		log.Println("Failed on call to SteamFriends() [pid is %p]\n", ptr)
		return ""
	}

	// Gets the AppId (previously set)
	var appId uint32
	fmt.Sscanf(os.Getenv("SteamAppId"), "%d", &appId)

	// Fails if it nil/null (0)
	if appId == 0 {
		log.Println("Failed to get AppID from env\n")
		return ""
	}

	// Gets the name (as a string)
	r1, _, _ := pISteamGetPersonaName.Call(ptr)

	// Copy it to a byte array
	buf := (*[unsafe.Sizeof(r1) - 1]byte)(unsafe.Pointer(r1))[:]

	// Strip the weird null pointers off the end.
	var i int
	for i = len(buf) - 1; i > 0; i-- {
		if buf[i] == 0x00 {
			break
		}
	}

	// Cast to string and then return
	return string(buf[0:i])
}
