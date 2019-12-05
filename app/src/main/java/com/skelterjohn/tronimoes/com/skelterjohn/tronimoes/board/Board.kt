package com.skelterjohn.tronimoes.board

data class Tile(val left: Face, val right: Face) {
    init {
        left.twin = right
        right.twin = left
    }
}

class Face(_pips: Int) {
    val pips = _pips

    var loc: V2? = null

    var twin: Face? = null
    var connections: MutableSet<Face> = mutableSetOf<Face>()

    var player: Player? = null

    var rank: Rank? = null

    fun open(): Boolean {

        if (rank == Rank.ROUND_LEADER) {
            return connections.size < 3
        }
        if (rank == Rank.LEADER) {
            return connections.isEmpty()
        }
        if (rank == Rank.LINE) {
            // Can play off either side of a double.
            if (pips == twin!!.pips) {
                return (connections + twin!!.connections).size == 1
            }
            // If this side is empty, the other must have a connection.
            return connections.isEmpty()
        }
        return false
    }
}

data class Player(val name: String) {
    var living = true
    var victims = mutableSetOf<Player>()
}

class Board(_width: Int, _height: Int) {
    val width = _width
    val height = _height

    // All the players who are chickenfooted.
    var chickenFeet = mutableSetOf<Player>()
    var players = mutableSetOf<Player>()

    var grid = Array<Face?>(width*height) { null }
    var tiles = mutableSetOf<Tile>()

    fun at(loc: V2): Face? {
        if (loc.x < 0 || loc.y < 0 || loc.x >= width || loc.y >= width) {
            return null
        }
        val i = loc.x + width*loc.y
        return grid[i]
    }
    fun put(face: Face, loc: V2) {
        val i = loc.x + width*loc.y
        grid[i] = face
    }

    fun openPips(player: Player, loc: V2): Int? {
        // No tile? Nothing is open.
        val f = at(loc) ?: return null

        // You can't connect to a start marker (only replace it).
        if (f.rank == Rank.START_MARKER) {
            return null
        }

        // If it's the round leader, this player can't already have a connected tile.
        if (f.rank == Rank.ROUND_LEADER) {
            for (connection in f.connections + f.twin!!.connections) {
                var cplayer = connection.player ?: continue
                if (cplayer == player) return null
            }
        }

        // If it's a leader, no one can have played off this tile.
        if (f.rank == Rank.LEADER) {
            if (!f.connections.isEmpty() || !f.twin!!.connections.isEmpty()) {
                return null
            }
        }

        // If not this player's tile, that tile must be a leader or round leaader, or that
        // player must be chickenfooted.
        if (player != f.player) {
            if (f.rank == Rank.LINE && f.player !in chickenFeet) {
                return null
            }
        }

        if (f.open()) {
            return f.pips
        }

        return null
    }

    fun placeFace(player: Player, face: Face, loc: V2, rank: Rank): Boolean {
        face.loc = loc
        face.rank = rank

        if (rank == Rank.LINE) {
            // Find what this tile links to. All valid parents will be used.
            for (adjLoc in loc.adjacent()) {
                if (face.pips == openPips(player, adjLoc)) {
                    var parent: Face = at(adjLoc) ?: continue
                    face.connections.add(parent)
                }
            }
            return !face.connections.isEmpty()
        }
        // If we didn't find any connections, this needs to be a leader.
        if (rank == Rank.ROUND_LEADER) {
            // It must be in the center of the board.
            return (loc == V2(width/2-1, height/2) || loc == V2(width/2, height/2))
        }
        if (rank == Rank.LEADER) {

        }
        if (rank == Rank.START_MARKER) {

        }
        return false
    }

    fun place(player: Player, tile: Tile, placement: Placement, rank: Rank): Boolean {
        // Ensure these locations are not occupied.
        if (at(placement.left) != null && at(placement.right) != null) {
            return false
        }

        // At least one of the faces has to be successfully placed.
        val placedLeft = placeFace(player, tile.left, placement.left, rank)
        val placedRight = placeFace(player, tile.right, placement.right, rank)
        if (!placedLeft && !placedRight) {
            return false
        }

        for (c in tile.left.connections) {
            c.connections.add(tile.left)
        }

        for (c in tile.right.connections) {
            c.connections.add(tile.right)
        }

        if (rank == Rank.LINE) {
            var connections = tile.left.connections + tile.right.connections
            // Only LINE tiles inherit the player.
            if (connections.size == 1) {
                var parent = connections.elementAt(0)
                if (parent.rank == Rank.LINE) {
                    tile.left.player = parent.player
                    tile.right.player = parent.player
                } else {
                    tile.left.player = player
                    tile.right.player = player
                }
            } else {
                // More than one connection is il ouroboros. Everyone involved dies.
                kill(player, player)
                for (face in connections) {
                    kill(player, face.player)
                }
            }
        }

        put(tile.left, placement.left)
        put(tile.right, placement.right)
        tiles.add(tile)

        // Identify the freshly killed.
        for (p in players) {
            if (!p.living) {
                continue
            }
            if (!canPlay(p)) {
                kill(player, p)
            }
        }

        return true
    }

    fun canPlay(player: Player): Boolean {
        return true
    }

    fun kill(victor: Player, victim: Player?) {
        victor.victims.add(victim ?: return)
        victim.living = false
    }
}

data class V2(val x: Int, val y: Int) {
    operator fun plus(o: V2): V2 {
        return V2(x + o.x, y + o.y)
    }
    operator fun minus(o: V2): V2 {
        return V2(x - o.x, y - o.y)
    }
    fun adjacent(): Set<V2> {
        var adj = mutableSetOf<V2>()
        adj.add(V2(x-1, y))
        adj.add(V2(x+1, y))
        adj.add(V2(x, y-1))
        adj.add(V2(x, y+1))
        return adj
    }
}

enum class Rank {
    LINE, LEADER, ROUND_LEADER, START_MARKER
}

data class Placement(val left: V2, val right: V2) {
    init {
        val delta = left - right
        assert(delta == V2(0, 1) || delta == V2(0, -1) || delta == V2(-1, 0) || delta == V2(1, 0))
    }
}
