const http = require('http'); // Use 'https' se sua URL for https

const API_CONFIG = {
    hostname: '192.168.2.220',
    port: 8080,
    path: '/orders',
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
};

const USERS = [
    "696028aff2ba343bf6310796",
    "696028aff2ba343bf6310796",
    "696028aff2ba343bf6310796"
];

const PRODUCTS = [
    "6960282df2ba343bf6310793",
    "6966aeaff2ba343bf631079b",
    "6966aefff2ba343bf631079c",
    "6966af29f2ba343bf631079d",
    "6966af4ff2ba343bf631079e"
]

const TOTAL_REQUESTS = 50;

function sendRequest(payload) {
    return new Promise((resolve, reject) => {
        const req = http.request(API_CONFIG, (res) => {
            let data = '';
            res.on('data', (chunk) => data += chunk);
            res.on('end', () => resolve(data));
        });

        req.on('error', (err) => reject(err));
        req.write(JSON.stringify(payload));
        req.end();
    });
}

async function startSimulation() {
    console.log(`🚀 Iniciando massa de dados nativa (${TOTAL_REQUESTS} pedidos)...`);

    for (let i = 1; i <= TOTAL_REQUESTS; i++) {
        const randomUser = USERS[Math.floor(Math.random() * USERS.length)];
        const randomProduct = PRODUCTS[Math.floor(Math.random() * PRODUCTS.length)];
        const randomQuantity = Math.floor(Math.random() * 5) + 1;

        const payload = {
            "user_id": randomUser,
            "items": [
                {
                    "product_id": randomProduct,
                    "quantity": randomQuantity
                }
            ],
            "card_token": "tok_visa"
        };

        try {
            const response = await sendRequest(payload);
            console.log(`✅ [${i}/${TOTAL_REQUESTS}] User: ${randomUser.slice(-4)} | Qtd: ${randomQuantity}`);
        } catch (err) {
            console.error(`❌ Erro no pedido ${i}:`, err.message);
        }

        // Intervalo aleatório entre 1s e 5s
        const waitTime = Math.floor(Math.random() * 4000) + 1000;
        console.log(`⏱️ Aguardando ${waitTime}ms...`);
        await new Promise(resolve => setTimeout(resolve, waitTime));
    }
}

startSimulation();
