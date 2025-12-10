const express = require("express");
const mysql = require("mysql2/promise");

const app = express();
const port = 3000;

async function getConnection() {
  return mysql.createConnection({
    host: "mysql",
    user: "root",
    password: "root",
    database: "fullcycle",
  });
}

app.get("/", async (req, res) => {
  const connection = await getConnection();

  await connection.execute(
    "INSERT INTO people (name) VALUES ('Daniel Cardoso')"
  );

  const [rows] = await connection.execute(
    "SELECT name FROM people ORDER BY id DESC LIMIT 1"
  );

  const name = rows[0].name;

  return res.send(`<h1>Full Cycle - ${name}</h1>`);
});

app.listen(port, () => {
  console.log("Rodando na porta " + port);
});
