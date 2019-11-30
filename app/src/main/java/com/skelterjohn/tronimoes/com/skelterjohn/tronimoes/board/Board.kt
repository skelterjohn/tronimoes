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
        val i = loc.x + width*loc.y
        return grid[i]
    }

    //
    //fun canPlace(tile: Tile, loc: V2): Array<V2> {
//
  //  }

    fun place(parent: Tile?, tile: Tile, placement: Placement): Boolean {
        grid[placement.left.x+placement.left.y*width] = tile
        grid[placement.right.x+placement.right.y*width] = tile
        tile.placement = placement
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
}

enum class Rank {
    LINE, LEADER, ROUND_LEADER
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

    // The tile this one led from, if any.
    var parent: Tile? = null
    // This tile's child, or children if it has rank ROUND_LEADER
    var children = mutableSetOf<Tile>()
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
