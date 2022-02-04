package main

import (
	"fmt"
	"sort"
)

type Row struct {
	servers []*Server
}

func algorithm(config Config, unMap unavailablesMap, initialServers []*Server) {

	servers := sortServers(initialServers)

	rows := make([]Row, config.rows)

	// Assign Server
	for sPos := 0; sPos < len(servers); sPos++ {
		server := servers[sPos]
		placeServer(config, unMap, server, sPos, rows)
	}

	// Assign pool
	currentPool := 0
	for rPos := 0; rPos < len(rows); rPos++ {
		for sPos := 0; sPos < len(rows[rPos].servers); sPos++ {
			rows[rPos].servers[sPos].assignedPool = currentPool % config.nPools
			currentPool++
		}
	}
}

func placeServer(config Config, unMap unavailablesMap, server *Server, sPos int, rows []Row) {
	for i := 0; i < config.rows; i++ {
		for j := 0; j < config.slots; j++ {
			if j+server.size > config.slots {
				break
			}
			canFit := true
			for k := j; k < j+server.size && k < config.slots; k++ {
				if ok := unMap[fmt.Sprintf("%d %d", i, k)]; ok {
					canFit = false
					break
				}
			}
			if !canFit {
				continue
			}

			server.assignedRow = i
			server.assignedSlot = j
			server.assigned = true

			// Slots unavailable
			if rows[i].servers == nil {
				rows[i].servers = make([]*Server, 0)
			}
			rows[i].servers = append(rows[i].servers, server)
			for k := j; k < j+server.size && k < config.slots; k++ {
				unMap[fmt.Sprintf("%d %d", i, k)] = true
			}
			return
		}
	}
}

func sortServers(initialServers []*Server) []*Server {
	sort.Slice(initialServers, func(i, j int) bool {
		a := initialServers[i]
		b := initialServers[j]

		if a.capacity/a.size > b.capacity/b.size {
			return true
		}
		return false
	})

	// part1 := servers[:len(servers)/2]
	// part2 := servers[len(servers)/2:]

	part1 := make([]*Server, len(initialServers)/2+1)
	part2 := make([]*Server, len(initialServers)/2)
	for i := 0; i < len(initialServers); i += 2 {
		part1[i/2] = initialServers[i]
		if (i / 2) < len(part2) {
			part2[(i / 2)] = initialServers[i+1]
		}
	}

	sort.Slice(part1, func(i, j int) bool {
		a := part1[i]
		b := part1[j]

		if a.capacity/a.size > b.capacity/b.size {
			return true
		}
		return false
	})
	sort.Slice(part2, func(i, j int) bool {
		a := part2[i]
		b := part2[j]

		if a.capacity/a.size < b.capacity/b.size {
			return true
		}
		return false
	})
	servers := append(part1, part2...)
	return servers
}
