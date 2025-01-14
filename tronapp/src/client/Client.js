class Client {
    constructor(baseURL, name, userInfo) {
        this.baseURL = baseURL;
        this.name = name;
        this.userInfo = userInfo;
		this.playerID = userInfo?.uid;
    }


    async JoinGame(code, name) {
		this.name = name;
        return this.put(`/game/${code}`, { name });
    }

    async GetGame(code, version) {
        return this.get(`/game/${code}?version=${version}`);
    }

    async StartRound(code) {
        return this.post(`/game/${code}/start`);
    }

    async LayTile(code, tile) {
        return this.post(`/game/${code}/tile`, tile);
    }

    async LaySpacer(code, spacer) {
        return this.post(`/game/${code}/spacer`, spacer);
    }

    async DrawTile(code) {
        return this.post(`/game/${code}/draw`, {});
    }

    async SetChickenFoot(code, url) {
        return this.post(`/game/${code}/foot`, { url: url });
    }

    async Pass(code, selected) {
        return this.post(`/game/${code}/pass`, selected);
    }

    async LeaveOrQuit(code) {
        return this.post(`/game/${code}/leave`, {});
    }

    async RegisterPlayerName(name) {
        return this.post(`/players`, { name: name });
    }

    async GetPlayer() {
        return this.get(`/players/${this.playerID}`);
    }

	async UpdatePlayer(playerInfo) {
		return this.put(`/players/${this.playerID}`, playerInfo);
	}

	async React(code, url) {
		return this.post(`/game/${code}/react`, { url: url });
	}

    async get(path) {
        return this.doRequest('GET', path);
    }


    async post(path, body) {
        return this.doRequest('POST', path, body);
    }


    async put(path, body) {
        return this.doRequest('PUT', path, body);
    }


    async doRequest(method, path, body = null) {
		let headers =  {
			'Content-Type': 'application/json',
			'Accept': 'application/json',
			'X-Player-Name': this.name,
		}
		if (this.userInfo) {
			headers['Authorization'] = `Bearer ${await this.userInfo.getIdToken(false)}`;
			headers['X-Player-Id'] = this.playerID;
		}
        const response = await fetch(`${this.baseURL}${path}`, {
            method,
            headers: headers,
            body: body ? JSON.stringify(body) : null,
        });

        const data = await response.json();

        if (!response.ok) {
            throw {
                status: response.status,
                data: data,
                message: `HTTP error! status: ${response.status}`
            };
        }

        return data;
    }
}


export default function clientFor(name, userInfo) {
	const baseURL = process.env.NEXT_PUBLIC_API_URL || 
                   `${window.location.protocol}//${window.location.hostname}:8080`;
    return new Client(baseURL, name, userInfo);
}
