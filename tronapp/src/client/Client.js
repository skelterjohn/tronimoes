class Client {
    constructor(baseURL, name, key) {
        this.baseURL = baseURL;
		this.name = name;
		this.key = key;
    }


	async JoinGame(code, name) {
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

	async DrawTile(code) {
		return this.post(`/game/${code}/draw`, {});
	}

	async Pass(code) {
		return this.post(`/game/${code}/pass`, {});
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
        const response = await fetch(`${this.baseURL}${path}`, {
            method,
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
				'X-Player-Name': this.name,
				'Authorization': `Bearer ${this.key}`,
            },
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


export default function clientFor(name, key) {
	return new Client(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080', name, key);
}
