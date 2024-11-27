class Client {
    constructor(baseURL, name, userid, key) {
        this.baseURL = baseURL;
        this.name = name;
        this.userid = userid;
        this.key = key;
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

    async GetPlayerName() {
        return this.get(`/players`);
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
		if (this.key !== undefined) {
			headers['Authorization'] = `Bearer ${this.key}`;
		}
		if (this.userid !== undefined) {
			headers['X-Player-Id'] = this.userid;
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


export default function clientFor(name, id, key) {
    return new Client(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080', name, id, key);
}
