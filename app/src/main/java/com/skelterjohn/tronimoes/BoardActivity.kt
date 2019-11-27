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

import com.skelterjohn.tronimoes.board.*
import kotlinx.android.synthetic.main.activity_board.*

class BoardView @JvmOverloads constructor(context: Context,
                                           attrs: AttributeSet? = null, defStyleAttr: Int = 0)
    : View(context, attrs, defStyleAttr) {

    public var board: Board? = null

    /// Position and scale factors for drawing.
    // The board coordinate that appears in the exact middle of the display. The round leader is
    // always at board coordinate 0,0-1,0.
    var centerX: Float = 0F
    var centerY: Float = 0F
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
        var pd = ProjectionDraw(canvas, G2(centerX, centerY), scaleFactor)
        for (t in board?.tiles ?: mutableSetOf<Tile>()) {
            pd.drawTile(t)
        }
    }


    data class G2(val gx: Float, val gy: Float) {
        constructor(v: V2) : this(v.x.toFloat(), v.y.toFloat())
        operator fun plus(o: G2): G2 {
            return G2(gx+o.gx, gy+o.gy)
        }
        fun turn(): G2 {
            return G2(gy, -gx)
        }
    }

    data class GRect(val tl: G2, val br: G2)

    class ProjectionDraw(canvas: Canvas, center: G2, scaleFactor: Float) {

        private val redPaint =
            Paint().apply {
                isAntiAlias = true
                color = Color.RED
                style = Paint.Style.STROKE
            }
        private val blackPaint =
            Paint().apply {
                isAntiAlias = true
                color = Color.BLACK
                style = Paint.Style.STROKE
            }
        var canvas: Canvas
        /// Position and scale factors for drawing.
        // The board coordinate that appears in the exact middle of the display. The round leader is
        // always at board coordinate 0,0-1,0.
        var center: G2
        // How many tile units fit between the center and the edge of the screen, in the smaller
        // dimension.
        var scaleFactor: Float
        var span: Int

        init {
            this.canvas = canvas
            this.center = center
            this.scaleFactor = scaleFactor
            span = canvas.width
            if (canvas.height < span) {
                span = canvas.height
            }
        }

        fun drawTile(tile: Tile) {
            val origin = G2(tile.origin ?: V2(0, 0))
            val delta = G2(tile.delta ?: V2(0, 0))

            var tl = origin
            var br = origin+delta+G2(1F, 1F)
            if (delta.gx < 0 || delta.gy < 0) {
                tl = origin+delta
                br = origin+G2(1F, 1F)
            }
            drawGRect(GRect(tl, br), blackPaint)

            drawPips(origin, delta, tile.left)
            drawPips(origin+delta, delta,tile.right)

            if (delta.gx == 0F) {
                drawGLine(G2(0.2F*tl.gx + 0.8F*br.gx,(tl.gy+br.gy)/2), G2(0.8F*tl.gx + 0.2F*br.gx,(tl.gy+br.gy)/2), redPaint)
            } else {
                drawGLine(G2((tl.gx+br.gx)/2, 0.2F*tl.gy + 0.8F*br.gy), G2((tl.gx+br.gx)/2, 0.8F*tl.gy + 0.2F*br.gy), redPaint)
            }
        }

        fun drawGRect(r: GRect, paint: Paint) {
            val tl = mapG2(r.tl)
            val br = mapG2(r.br)
            canvas.drawRect(tl.dx, tl.dy, br.dx, br.dy, paint)
        }

        fun drawGLine(start: G2, end: G2, paint: Paint) {
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

        val pipLocatoinSet = arrayOf<Array<G2>>(
            arrayOf<G2>(), // 0
            arrayOf<G2>(G2(0F, 0F)), // 1
            arrayOf<G2>(G2(-0.2F, -0.2F),
                        G2(0.2F, 0.2F)), // 2
            arrayOf<G2>(G2(-0.2F, -0.2F),
                        G2(0F, 0F),
                        G2(0.2F, 0.2F)), // 3
            arrayOf<G2>(G2(-0.2F, -0.2F),
                        G2(0.2F, -0.2F),
                        G2(-0.2F, 0.2F),
                        G2(0.2F, 0.2F)), // 4
            arrayOf<G2>(G2(-0.2F, -0.2F),
                        G2(0.2F, -0.2F),
                        G2(0F, 0F),
                        G2(-0.2F, 0.2F),
                        G2(0.2F, 0.2F)), // 5
            arrayOf<G2>(G2(-0.25F, -0.2F),
                        G2(0.25F, -0.2F),
                        G2(0F, -0.2F),
                        G2(0F, 0.2F),
                        G2(-0.25F, 0.2F),
                        G2(0.25F, 0.2F)) // 6
        )

        fun drawPips(loc: G2, delta: G2, pips: Int) {
            val pipLocs = pipLocatoinSet[pips]
            for (i in 0..pipLocs.size-1) {
                var pl = pipLocs[i]
                when(delta) {
                    G2(0F, 1F) -> pl = pl.turn()
                    G2(1F, 0F) -> pl = pl.turn().turn()
                    G2(0F, -1F) -> pl = pl.turn().turn().turn()
                }
                drawPip(loc+G2(0.5F, 0.5F)+pl, redPaint)
            }
        }

        fun drawPip(loc: G2, paint: Paint) {
            val pc = mapG2(loc)
            canvas.drawCircle(pc.dx, pc.dy, scale(0.1F), redPaint)
        }
    }
}

class BoardActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_board)
        var b = Board(50, 50)
        board_view.board = b

        var p = Pile(6)
        var t1 = p.draw("john")
        t1.origin = V2(0, 0)
        t1.delta = V2(1, 0)
        var t2 = p.draw("stef")
        t2.origin= V2(0, 1)
        t2.delta = V2(0, 1)
        var t3 = p.draw("john")
        t3.origin = V2(2, 0)
        t3.delta = V2(1, 0)

        b.tiles.add(t1)
        b.tiles.add(t2)
        b.tiles.add(t3)
    }

}
