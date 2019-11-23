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

    // Called when the view should render its content.
    override fun onDraw(canvas: Canvas?) {
        super.onDraw(canvas)
        if (canvas == null) {
            return
        }
        // DRAW STUFF HERE
        var color = redPaint
        if (board != null) {
            color = blackPaint
        }

        for (t in board?.tiles ?: mutableSetOf<Tile>()) {
            drawTile(canvas, t)
        }
    }

    fun drawTile(canvas: Canvas, tile: Tile) {

        var span = canvas.width
        if (canvas.height < span) {
            span = canvas.height
        }
        val origin: V2 = tile.origin?: V2(0, 0)
        val delta: V2 = tile.delta?: V2(0, 0)
        var leftX = origin.x.toFloat()
        var leftY = origin.y.toFloat()
        var leftMinX = mapX(leftX-0.5F, centerX, canvas.width/2, span)
        var leftMaxX = mapX(leftX+0.5F, centerX, canvas.width/2, span)
        var leftMinY = mapY(leftY-0.5F, centerY, canvas.height/2, span)
        var leftMaxY = mapY(leftY+0.5F, centerY, canvas.height/2, span)
        canvas.drawRect(RectF(leftMinX, leftMaxY, leftMaxX, leftMinY), blackPaint)
        var rightX = (origin+delta).x.toFloat()
        var rightY = (origin+delta).y.toFloat()
        var rightMinX = mapX(rightX-0.5F, centerX, canvas.width/2, span)
        var rightMaxX = mapX(rightX+0.5F, centerX, canvas.width/2, span)
        var rightMinY = mapY(rightY-0.5F, centerY, canvas.height/2, span)
        var rightMaxY = mapY(rightY+0.5F, centerY, canvas.height/2, span)
        canvas.drawRect(RectF(rightMinX, rightMaxY, rightMaxX, rightMinY), blackPaint)
    }

    fun mapX(gx: Float, gcx: Float, dcx: Int, span: Int): Float {
        // distance from the center to the edge.
        var offset = gx - gcx
        var dx = dcx + (offset * (span/2) / scaleFactor)
        return dx
    }
    fun mapY(gy: Float, gcy: Float, dcy: Int, span: Int): Float {
        // distance from the center to the edge.
        var offset = gy - gcy
        var dy = dcy - (offset * (span/2) / scaleFactor)
        return dy
    }
}

class BoardActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_board)
        var b = Board(50, 50)
        board_view.board = b

        var p = Pile(5)
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
