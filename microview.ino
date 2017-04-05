#include <MicroView.h>

void setup() {
  uView.begin();
  uView.clear(PAGE);

  Serial.begin(115200);
  Serial.print("MicroView");
}

void loop() {
  uView.checkComm();
}
