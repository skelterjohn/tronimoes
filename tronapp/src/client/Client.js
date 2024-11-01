class Client {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }


	async join(name, code) {
		return this.put(`/game/${code}`, { name });
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


// Export a singleton instance
const client = new Client(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080');
export default client;
