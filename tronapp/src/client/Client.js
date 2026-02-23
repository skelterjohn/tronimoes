class Client {
    constructor(baseURL, name, userInfo) {
        this.baseURL = baseURL;
        this.name = name;
        this.userInfo = userInfo;
		this.playerID = userInfo?.uid;
    }


    async JoinGame(code, name, options) {
		this.name = name;
        return this.put(`/game/${code}`, options);
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

	async UpdatePlayerConfig(config) {
		return this.put(`/players/${this.playerID}/config`, config);
	}

	async React(code, url) {
		return this.post(`/game/${code}/react`, { url: url });
	}

	async ReportIssue(code, summary, whatHappened, whatShouldHappen, errorMessage) {
		return this.post(`/game/${code}/report`, {
			summary: summary,
			what_happened: whatHappened,
			what_should_happen: whatShouldHappen,
			error_message: errorMessage,
		});
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
		let response;
		try {
			response = await fetch(`${this.baseURL}${path}`, {
				method,
				headers: headers,
				body: body ? JSON.stringify(body) : null,
			});
		} catch (err) {
			const isCorsOrNetwork = err instanceof TypeError ||
				(err?.message && (
					String(err.message).toLowerCase().includes('failed to fetch') ||
					String(err.message).toLowerCase().includes('network error')
				));
			if (isCorsOrNetwork) {
				throw {
					status: 0,
					data: { error: `Request problem: ${err?.message}`},
					message: err?.message || 'Failed to fetch'
				};
			}
			throw err;
		}
		let data;
		try {
			data = await response.json();
		} catch {
			data = { error: 'Invalid JSON response' };
		}

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
