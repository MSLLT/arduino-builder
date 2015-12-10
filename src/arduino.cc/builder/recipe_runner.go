/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/utils"
	"os"
	"sort"
	"strings"
)

type RecipeByPrefixSuffixRunner struct {
	Prefix string
	Suffix string
}

func (s *RecipeByPrefixSuffixRunner) Run(context map[string]interface{}) error {
	logger := context[constants.CTX_LOGGER].(i18n.Logger)
	if utils.DebugLevel(context) >= 10 {
		logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_LOOKING_FOR_RECIPES, s.Prefix, s.Suffix)
	}

	buildProperties := utils.GetMapStringStringOrDefault(context, constants.CTX_BUILD_PROPERTIES)
	verbose := context[constants.CTX_VERBOSE].(bool)

	recipes := findRecipes(buildProperties, s.Prefix, s.Suffix)

	properties := utils.MergeMapsOfStrings(make(map[string]string), buildProperties)
	for _, recipe := range recipes {
		if utils.DebugLevel(context) >= 10 {
			logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_RUNNING_RECIPE, recipe)
		}
		_, err := builder_utils.ExecRecipe(properties, recipe, false, verbose, verbose, logger)
		if err != nil {
			return utils.WrapError(err)
		}
	}

	return nil

}

func findRecipes(buildProperties map[string]string, patternPrefix string, patternSuffix string) []string {
	var recipes []string
	for key, _ := range buildProperties {
		if strings.HasPrefix(key, patternPrefix) && strings.HasSuffix(key, patternSuffix) && buildProperties[key] != constants.EMPTY_STRING {
			recipes = append(recipes, key)
		}
	}

	sort.Strings(recipes)

	return recipes
}
