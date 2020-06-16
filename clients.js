async function getResults(accessToken) {
    const requestInit = {
        headers: {
            Authorization: accessToken
        },
        mode: 'cors'
    }
    console.log(requestInit)
    const results = await fetch('http://localhost:8081/results', requestInit)
    const resultsJson = results.json()
    return resultsJson
}

async function getTransports(accessToken) {
    const requestInit = {
        headers: {
            Authorization: accessToken
        },
        mode: 'cors'
    }
    const transports = await fetch('http://localhost:8081/transports', requestInit)
    const transportsJson = transports.json()
    return transportsJson
}

async function getRoutes(accessToken, busId) {
    if (!busId || !accessToken) return;
    const requestInit = {
        headers: {
            Authorization: accessToken
        },
        mode: 'cors'
    }
    const transports = await fetch(`http://localhost:8081/routes/${busId}`, requestInit)
    const transportsJson = transports.json()
    return transportsJson
}

async function getUsers(accessToken) {
    const requestInit = {
        headers: {
            Authorization: accessToken
        },
        mode: 'cors'
    }
    const users = await fetch('http://localhost:8081/admin/users', requestInit)
    const usersJson = users.json()
    return usersJson
}

async function login(username, password) {
    const requestinit = {
        headers: {
            "Content-Type": "application/json",
        },
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify({
            Username: username,
            Password: password
        })
    }
    return fetch('http://localhost:8081/login', requestinit)
}

async function AddResult(accessToken, checkupPayload) {
    const requestinit = {
        headers: {
            Authorization: accessToken,
            "Content-Type": "application/json",
        },
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(checkupPayload)
    }
    return fetch('http://localhost:8081/results', requestinit)
}

async function registerUser(accessToken, registerPayload) {
    const requestinit = {
        headers: {
            Authorization: accessToken,
            "Content-Type": "application/json",
        },
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(registerPayload)
    }
    return fetch('http://localhost:8081/admin/users', requestinit)
}

export {
    getResults,
    getTransports,
    getUsers,
    login,
    AddResult,
    registerUser,
    getRoutes
}