package com.skelterjohn.tronimoes.com.skelterjohn.tronimoes.board

import kotlin.random.Random

class Board(_width: Int, _height: Int) {
    val width = _width
    val height = _height

    // All tiles played on this board.
    var tiles = mutableSetOf<Tile>()
    // All tiles that are currently available leaders.
    var leaders = mutableSetOf<Tile>()
}

class V2(_x: Int, _y: Int) {
    val x = _x
    val y = _y
}

class Tile(_left: Int, _right: Int) {
    // The number of pips on the left side.
    val left = _left
    // The number of pips on the right side.
    val right = _right

    // The player who placed this tile.
    var player: String? = null

    // The position of the left and right sides.
    var leftCoord: V2? = null
    var rightCoord: V2? = null

    // The tile this one led from, if any.
    var parent: Tile? = null
}

class Pile(maxPips:Int) {
    var tiles: MutableList<Tile> = mutableListOf<Tile>()

    init {
        for (left in 0..maxPips) {
            for (right in 0..maxPips) {
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
