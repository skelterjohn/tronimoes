package com.skelterjohn.tronimoes.board

import kotlin.random.Random

class Board(_width: Int, _height: Int) {
    val width = _width
    val height = _height

    // All tiles played on this board.
    var tiles = mutableSetOf<Tile>()
    // All tiles that are currently available leaders.
    var leaders = mutableSetOf<Tile>()
    // All the players who are chickenfooted.
    var chickenFeet = mutableSetOf<String>()

    var grid = Array<Tile?>(width*height) { i -> null }

    fun tile(loc: V2): Tile? {
        if (loc.x < 0 || loc.y < 0 || loc.x >= width || loc.y >= width) {
            return null
        }
        val i = loc.x + width*loc.y
        return grid[i]
    }

    fun openPips(player: String, loc: V2): Int? {
        // No tile? Nothing is open.
        val t = tile(loc) ?: return null
        val tplayer = t.player ?: ""
        val tplacement = t.placement ?: return null

        // If it's the round leader, this player can't already have a connected tile.
        if (t.rank == Rank.ROUND_LEADER) {
            for (connection in t.leftConnections + t.rightConnections) {
                var cplayer = connection?.player ?: continue
                if (cplayer == player) return null
            }
        }

        // If it's a leader, no one can have played off this tile.
        if (t.rank == Rank.LEADER) {
            if (!t.leftConnections.isEmpty() || !t.rightConnections.isEmpty()) {
                return null
            }
        }

        // If you're chickenfooted it needs to be your own tile or the round leader.
        if (player in chickenFeet && t.rank != Rank.ROUND_LEADER && player != t.player) return null

        // If it's another player's tile, that tile must be a leader or round leaader, or that
        // player must be chickenfooted.
        if (tplayer != "" && player != tplayer) {
            if (tplayer !in chickenFeet) {
                if (t.rank != Rank.ROUND_LEADER && t.rank != Rank.LEADER) {
                    return null
                }
            }
        }

        // Make sure the tile is open at loc.
        val ol = t.openLeft()
        val or = t.openRight()
        if (!ol && !or) return null
        if (loc == tplacement.left && !ol) return null
        if (loc == tplacement.right && !or) return null

        if (loc == tplacement.left) return t.left
        if (loc == tplacement.right) return t.right

        return null
    }

    // Connect a tile to a parent.
    fun connect(tile: Tile, tileLoc: V2, parent: Tile, parentLoc: V2) {
        if (tile.placement!!.left == tileLoc) {
            tile.leftConnections.add(parent)
        }
        if (tile.placement!!.right == tileLoc) {
            tile.rightConnections.add(parent)
        }
        if (parent.placement!!.left == parentLoc) {
            parent.leftConnections.add(tile)
        }
        if (parent.placement!!.right == parentLoc) {
            parent.rightConnections.add(tile)
        }
    }

    fun place(tile: Tile, placement: Placement): Boolean {
        tile.placement = placement
        // Ensure these locations are not occupied.
        if (tile(placement.left) != null && tile(placement.right) != null) {
            return false
        }
        val player = tile.player ?: return false
        if (tile.rank == Rank.LINE) {
            // Find what this tile links to. All valid parents will be used.
            for (loc in placement.left.adjacent()) {
                if (tile.left == openPips(player, loc)) {
                    var parent: Tile = tile(loc) ?: continue
                    connect(tile, placement.left, parent, loc)
                }
            }
            for (loc in placement.right.adjacent()) {
                if (tile.right == openPips(player, loc)) {
                    var parent: Tile = tile(loc) ?: continue
                    connect(tile, placement.right, parent, loc)
                }
            }
            if (tile.leftConnections.isEmpty() && tile.rightConnections.isEmpty()) {
                return false
            }
        }
        // If we didn't find any connections, this needs to be a leader.
        if (tile.rank == Rank.ROUND_LEADER) {
            if (placement != Placement(V2(width/2-1, height/2), V2(width/2, height/2))) {
                return false
            }
        }
        if (tile.rank == Rank.LEADER) {

        }
        if (tile.rank == Rank.START_MARKER) {

        }
        grid[placement.left.x+placement.left.y*width] = tile
        grid[placement.right.x+placement.right.y*width] = tile
        tiles.add(tile)
        return true
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

data class Placement(val left: V2, val right: V2)

class Tile(_left: Int, _right: Int) {
    // The number of pips on the left side.
    val left = _left
    // The number of pips on the right side.
    val right = _right
    // The rank for this tile.
    var rank: Rank? = null

    // For LINE tiles, the player who owns it.
    var player: String? = null

    // The position of the left and right sides.
    var placement: Placement? = null

    // Which sides of this tile are open for children.
    fun open(connections: Set<Tile>, otherSideConnections: Set<Tile>): Boolean {
        if (rank == Rank.ROUND_LEADER) {
            return connections.size < 3
        }
        if (rank == Rank.LEADER) {
            return connections.isEmpty()
        }
        if (rank == Rank.LINE) {
            // Can play off either side of a double.
            if (left == right) {
                return connections.isEmpty() || otherSideConnections.isEmpty()
            }
            // If this side is empty, the other must have a connection.
            return connections.isEmpty()
        }
        return false
    }
    fun openLeft(): Boolean {
        return open(leftConnections, rightConnections)
    }
    fun openRight(): Boolean {
        return open(rightConnections, leftConnections)
    }

    // The tiles connected to the left.
    var leftConnections = mutableSetOf<Tile>()
    // The tiles connected to the right.
    var rightConnections = mutableSetOf<Tile>()
}

class Pile(maxPips:Int) {
    var tiles: MutableList<Tile> = mutableListOf<Tile>()

    init {
        for (left in 0..maxPips) {
            for (right in left..maxPips) {
                tiles.add(Tile(left, right))
            }
        }
    }

    fun draw(player: String): Tile {
        var t = tiles.removeAt(Random.nextInt(0, tiles.size))
        t.player = player
        return t
    }
}
