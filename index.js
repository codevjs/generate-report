const express = require('express');
const dotenv = require('dotenv');
const reportRouter = require('./routes/report');

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

app.use('/report', reportRouter);

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
