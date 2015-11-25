int BASE_PIN = 2;

void setup() {
  Serial.begin(9600);
  for (int i = 0; i < 4; i++) {
    pinMode(BASE_PIN + i, OUTPUT);
  }

  // Apaga todos os LEDs
  for (int i = 0; i < 4; i++) {
    digitalWrite(BASE_PIN + i, LOW);
  }
}

void loop() {
  int button, clicks;

  // Aguarda o recebimento de dados
  if (Serial.available() >= 2) {
    // Lê mensagem enviada
    button = Serial.read();
    clicks = Serial.read();

    // Informa qual botao foi pressionado através do LED 9
    blink(BASE_PIN, button+1);
    // Quantidade de vezes que o botão foi pressionado
    blink(BASE_PIN+button+1, clicks);

    // Informa que mensagem foi processada corretamente
    Serial.print(button, DEC);
    Serial.print(" ");
    Serial.println(clicks, DEC);
  }
}

void blink(int ledNumber, int times) {
  for (int i = 0; i < times; i++) {
    digitalWrite(ledNumber, HIGH);
    delay(800);
    digitalWrite(ledNumber, LOW);
    delay(300);
  }
}
