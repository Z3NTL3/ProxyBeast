package builder

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

import (
	"Z3NTL3/proxy-checker/fancy"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// .-,--.                   ,-,---.             .
//  '|__/ ,-. ,-. . , . .    '|___/ ,-. ,-. ,-. |-
//  ,|    |   | |  X  | |    ,|   \ |-' ,-| `-. |
//  `'    '   `-' ' ` `-|   `-^---' `-' `-^ `-' `'
//                     /|
//                    `-'

func randomRGB() (string){
	return fancy.Palettes[rand.Intn(len(fancy.Palettes))]
}

func Logo(){
	rand.Seed(time.Now().Unix())

	var logo []string
	logo = append(logo, fmt.Sprintf("%s%s%s%s%s%s%s", fancy.Bold, randomRGB() , "\033[1m .-,","--.                   ", randomRGB(),",-,-","--.             ."))
	logo = append(logo, fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s", fancy.Bold, randomRGB() ,"  '|__/",randomRGB() ," ,-.",randomRGB()," ,-. ",randomRGB(),". , . ",randomRGB(),".    '|___",randomRGB(),"/ ,-. ",randomRGB(),",-. ,",randomRGB(),"-. |-"))
	logo = append(logo, fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s",fancy.Bold,randomRGB(), "  ,|   ",randomRGB()," |   |",randomRGB()," |  X  ",randomRGB(),"| |    ",randomRGB(),",|   \\ ",randomRGB(),"|-' ,",randomRGB(),"-| `",randomRGB(),"-. |"))
	logo = append(logo, fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s", fancy.Bold, randomRGB(),"  `'    ",randomRGB(),"'   `-",randomRGB(),"' ' ` ",randomRGB(),"`-|   ",randomRGB(),"`-^---",randomRGB(),"' `-' `-",randomRGB(),"^ `-' `'"))
	logo = append(logo, fmt.Sprintf("%s%s%s", fancy.Bold, randomRGB(), "                     /|"))
	logo = append(logo, fmt.Sprintf("%s%s%s",fancy.Bold, randomRGB(),"                    `-'\n"))
	logo = append(logo, fmt.Sprintf("\t%s%s%sTool by: @%sz3ntl3\033[0m", fancy.Bold, randomRGB(),randomRGB(), randomRGB()))
	logo = append(logo, fmt.Sprintf("\t%s%sStudios: %s\x1b]8;;/\ahttps://pix4.dev\x1b]8;;\a \033[0m\r\n", fancy.Bold, randomRGB(), randomRGB()))

	fmt.Println(strings.Join(logo, fmt.Sprintf("%s%s", fancy.Endline, fancy.Reset)))
}