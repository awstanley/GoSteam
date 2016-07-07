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
)

// InitSteamApp attempts to start a new SteamApp (setting it as active until we leave)
func InitSteamApp(AppID int64) (err error) {
	// Set the AppID
	os.Setenv("SteamAppId", fmt.Sprintf("%d", AppID))

	if !steamApiInit() {
		return fmt.Errorf("steam reported failure to launch")
	}

	log.Println("Successfully got application.")

	return nil
}
