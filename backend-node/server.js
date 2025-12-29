const express = require('express');
const cors = require('cors');
const app = express();
const port = 8080;

app.use(cors());
app.use(express.json());

const assets = [
    {
        symbol: "BTC",
        mark_price: 62000,
        contract_value: 0.001,
        allowed_leverage: [5, 10, 20, 50, 100]
    },
    {
        symbol: "ETH",
        mark_price: 3200,
        contract_value: 0.01,
        allowed_leverage: [5, 10, 25, 50]
    }
];

app.get('/config/assets', (req, res) => {
    res.json({ assets });
});

app.post('/margin/validate', (req, res) => {
    const { asset, order_size, side, leverage, margin_client } = req.body;

    const selectedAsset = assets.find(a => a.symbol === asset);
    if (!selectedAsset) {
        return res.status(400).json({ status: 'error', message: 'Asset not found', margin_required: 0 });
    }

    if (!selectedAsset.allowed_leverage.includes(leverage)) {
        return res.status(400).json({ status: 'error', message: 'Invalid leverage', margin_required: 0 });
    }

    let marginRequired = (selectedAsset.mark_price * order_size * selectedAsset.contract_value) / leverage;

    marginRequired = Math.round(marginRequired * 100) / 100;

    if (margin_client < marginRequired) {
        return res.status(400).json({
            status: 'error',
            message: 'Insufficient margin submitted',
            margin_required: marginRequired
        });
    }

    res.json({
        status: 'ok',
        margin_required: marginRequired
    });
});

app.listen(port, () => {
    console.log(`Node.js backend server listening on port ${port}`);
});
