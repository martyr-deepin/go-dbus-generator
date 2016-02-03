/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import (
	"strings"
)

func normalizeMethodName(name string) string {
	words := strings.Split(name, "_")
	normalized := ""
	for _, w := range words {
		normalized += upper(w)
	}
	return upper(normalized)
}
