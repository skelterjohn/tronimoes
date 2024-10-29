import Board from '../board/Board';
import Hand from './Hand';

function Game({}) {
	return <div className="">
		<table>
			<tbody>
				<tr></tr>
				<tr>
					<td>
						<Board width={10} height={11} />
					</td>
				</tr>
				<tr>
					<td>
						<Hand/>
					</td>
				</tr>
			</tbody>
		</table>
	</div>;
}

export default Game;