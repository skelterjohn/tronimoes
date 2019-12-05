package com.skelterjohn.tronimoes

import android.content.Context
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import android.graphics.Canvas
import android.graphics.Color
import android.graphics.Paint
import android.graphics.RectF
import android.util.AttributeSet
import android.util.Log
import android.view.View
import kotlin.math.min
import kotlin.math.max

import com.skelterjohn.tronimoes.board.*
import kotlinx.android.synthetic.main.activity_board.*

class HandView @JvmOverloads constructor(context: Context,
                                          attrs: AttributeSet? = null, defStyleAttr: Int = 0)
    : View(context, attrs, defStyleAttr) {

    var board: Board? = null

    // Called when the view should render its content.
    override fun onDraw(canvas: Canvas?) {
        super.onDraw(canvas)
        if (canvas == null) {
            return
        }
    }
}

class BoardView @JvmOverloads constructor(context: Context,
                                           attrs: AttributeSet? = null, defStyleAttr: Int = 0)
    : View(context, attrs, defStyleAttr) {

    var board: Board? = null
    var shades = mutableMapOf<Player, Int>()
    var turn: Int = 0
    var players = arrayOf<Player>()

    /// Position and scale factors for drawing.
    // The board coordinate that appears in the exact middle of the display. The round leader is
    // always at board coordinate 0,0-1,0.
    var center: G2 = G2(0F, 0F)
    // How many tile units fit between the center and the edge of the screen, in the smaller
    // dimension.
    var scaleFactor = 5F

    // Called when the view should render its content.
    override fun onDraw(canvas: Canvas?) {
        super.onDraw(canvas)
        if (canvas == null) {
            return
        }
        // DRAW STUFF HERE

        canvas.drawRect(RectF(0F, 0F, this.width.toFloat(), this.height.toFloat()),
            Paint().apply {
                isAntiAlias = true
                color = Color.parseColor("#CCFFFF")
                style = Paint.Style.FILL
            })

        var pd = ProjectionDraw(canvas, this)
        for (t in board?.tiles ?: mutableSetOf<Tile>()) {
            pd.drawTile(t)
        }
    }

}

data class G2(val gx: Float, val gy: Float) {
    constructor(v: V2) : this(v.x.toFloat(), v.y.toFloat())
    operator fun plus(o: G2): G2 {
        return G2(gx+o.gx, gy+o.gy)
    }
    operator fun minus(o: G2): G2 {
        return G2(gx-o.gx, gy-o.gy)
    }
    fun turn(): G2 {
        return G2(gy, -gx)
    }
}

data class GRect(val tl: G2, val br: G2)

class ProjectionDraw(canvas: Canvas, boardView: BoardView) {
    var canvas: Canvas = canvas
    /// Position and scale factors for drawing.
    // The board coordinate that appears in the exact middle of the display.
    var center: G2 = boardView.center
    // How many tile units fit between the center and the edge of the screen, in the smaller
    // dimension.
    var scaleFactor: Float = boardView.scaleFactor
    var span: Int

    var playerShades: Map<Player,Int> = boardView.shades

    init {
        span = canvas.width
        if (canvas.height < span) {
            span = canvas.height
        }
    }

    fun drawTile(tile: Tile) {
        val left = G2(tile.left.loc ?: return)
        val right = G2(tile.right.loc ?: return)

        val tmin = G2(min(left.gx, right.gx), min(left.gy, right.gy))
        val tmax = G2(max(left.gx, right.gx), max(left.gy, right.gy)) + G2(1F, 1F)

        var tileShade: Int = playerShades[tile.left.player] ?: Color.parseColor("#EEEEEE")
        rect(GRect(tmin, tmax), Paint().apply {
            isAntiAlias = true
            color = tileShade
            style = Paint.Style.FILL
        })
        line(tmin, G2(tmax.gx, tmin.gy), Paint().apply {
            isAntiAlias = true
            color = Color.GRAY
            style = Paint.Style.STROKE
            strokeWidth = 3F
        })
        line(G2(tmax.gx, tmin.gy), tmax, Paint().apply {
            isAntiAlias = true
            color = Color.GRAY
            style = Paint.Style.STROKE
            strokeWidth = 3F
        })

        pips(left, right-left, tile.left.pips)
        pips(right, left-right, tile.right.pips)

        if (left.gx == right.gx) {
            val midy = (tmin.gy + tmax.gy) / 2
            val minx = 0.8F*tmin.gx + 0.2F*tmax.gx
            val maxx = 0.2F*tmin.gx + 0.8F*tmax.gx
            line(G2(minx, midy), G2(maxx, midy), Paint().apply {
                isAntiAlias = true
                color = Color.GRAY
                style = Paint.Style.STROKE
            })
        } else {
            val midx = (tmin.gx + tmax.gx) / 2
            val miny = 0.8F*tmin.gy + 0.2F*tmax.gy
            val maxy = 0.2F*tmin.gy + 0.8F*tmax.gy
            line(G2(midx, maxy), G2(midx, miny), Paint().apply {
                isAntiAlias = true
                color = Color.GRAY
                style = Paint.Style.STROKE
            })
        }
    }

    fun rect(r: GRect, paint: Paint) {
        val tl = mapG2(r.tl)
        val br = mapG2(r.br)
        canvas.drawRect(tl.dx, tl.dy, br.dx, br.dy, paint)
    }

    fun line(start: G2, end: G2, paint: Paint) {
        val dStart = mapG2(start)
        val dEnd = mapG2(end)
        canvas.drawLine(dStart.dx, dStart.dy, dEnd.dx, dEnd.dy, paint)
    }

    data class D2(val dx: Float, val dy: Float)

    fun mapG2(g: G2): D2 {
        return D2(
            mapX(g.gx, center.gx, canvas.width/2),
            mapY(g.gy, center.gy, canvas.height/2))
    }

    fun mapX(gx: Float, gcx: Float, dcx: Int): Float {
        // distance from the center to the edge.
        var offset = gx - gcx
        var dx = dcx + (offset * (span/2) / scaleFactor)
        return dx
    }
    fun mapY(gy: Float, gcy: Float, dcy: Int): Float {
        // distance from the center to the edge.
        var offset = gy - gcy
        var dy = dcy - (offset * (span/2) / scaleFactor)
        return dy
    }
    fun scale(distance: Float): Float {
        return distance * (span/2) / scaleFactor
    }

    data class PipsDesc(val offsets: Array<G2>, val color: Int)
    val pipLocatoinSet = arrayOf<PipsDesc>(
        PipsDesc(arrayOf<G2>(), Color.BLACK), // 0
        PipsDesc(arrayOf<G2>(G2(0F, 0F)), Color.BLUE), // 1
        PipsDesc(arrayOf<G2>(G2(-0.2F, -0.2F),
                    G2(0.2F, 0.2F)), Color.GREEN), // 2
        PipsDesc(arrayOf<G2>(G2(-0.2F, -0.2F),
                    G2(0F, 0F),
                    G2(0.2F, 0.2F)), Color.YELLOW), // 3
        PipsDesc(arrayOf<G2>(G2(-0.2F, -0.2F),
                    G2(0.2F, -0.2F),
                    G2(-0.2F, 0.2F),
                    G2(0.2F, 0.2F)), Color.CYAN), // 4
        PipsDesc(arrayOf<G2>(G2(-0.2F, -0.2F),
                    G2(0.2F, -0.2F),
                    G2(0F, 0F),
                    G2(-0.2F, 0.2F),
                    G2(0.2F, 0.2F)), Color.GRAY), // 5
        PipsDesc(arrayOf<G2>(G2(-0.25F, -0.2F),
                    G2(0.25F, -0.2F),
                    G2(0F, -0.2F),
                    G2(0F, 0.2F),
                    G2(-0.25F, 0.2F),
                    G2(0.25F, 0.2F)), Color.MAGENTA) // 6
    )

    fun pips(loc: G2, delta: G2, pips: Int) {
        val pipsDesc = pipLocatoinSet[pips]
        val pipLocs = pipsDesc.offsets
        for (i in 0..pipLocs.size-1) {
            var pl = pipLocs[i]
            when(delta) {
                G2(0F, 1F) -> pl = pl.turn()
                G2(1F, 0F) -> pl = pl.turn().turn()
                G2(0F, -1F) -> pl = pl.turn().turn().turn()
            }
            val pc = mapG2(loc+G2(0.5F, 0.5F)+pl)
            canvas.drawCircle(pc.dx, pc.dy, scale(0.1F), Paint().apply {
                isAntiAlias = true
                color = pipsDesc.color
                style = Paint.Style.FILL
            })
            canvas.drawCircle(pc.dx, pc.dy, scale(0.1F), Paint().apply {
                isAntiAlias = true
                color = Color.BLACK
                style = Paint.Style.STROKE
            })
        }
    }
}

class BoardActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_board)
        var b = Board(10, 11)

        board_view.board = b
        board_view.center = G2(5F, 5F)

        hand_view.board = b

        var john = Player("john")
        var stef = Player("stef")

        b.players.add(john)
        b.players.add(stef)

        for (i in 0..5) {
            john.hand.add(b.pile.draw())
            stef.hand.add(b.pile.draw())
        }

        board_view.shades[john] = Color.parseColor("#FFAAAA")
        board_view.shades[stef] = Color.parseColor("#AAFFAA")

        var t1 = Tile(Face(6), Face(6))
        var t1placed = b.place(john, t1, Placement(V2(4, 5), V2(5, 5)), Rank.ROUND_LEADER)
        var t2 = Tile(Face(6), Face(1))
        var t2placed = b.place(stef, t2, Placement(V2(6, 5), V2(7, 5)), Rank.LINE)
        var t3 = Tile(Face(6), Face(3))
        var t3placed = b.place(john, t3, Placement(V2(4, 6), V2(3, 6)), Rank.LINE)
        Log.i("onCreate", "${t1placed} ${t2placed} ${t3placed}")
    }

}
