import {Row, Col} from 'react-bootstrap';

import Board from '../board/Board';
import Hand from './Hand';

function Game({}) {
	return <div className="">
		<Row className="flex justify-center items-center">
			<Col>
				<Hand color="red" hidden={true} dead={true} />
			</Col>
			<Col>
				<Hand color="blue" hidden={true} dead={false} />
			</Col>
		</Row>
		<Row>
			<Board width={10} height={11} />
		</Row>
		<Row>
			<Hand color="green"/>

		</Row>
	</div>;
}

export default Game;