#include <Wire.h>
#include <WiFi.h>
#include <SPI.h>
#include <Adafruit_BMP280.h>
#include <Redis.h>
#include <ArduinoJson.h>
#include <Preferences.h>


#define REDIS_ADDR "0.0.0.0"
#define REDIS_PORT 6379
#define REDIS_PASSWORD "password"
#define REDIS_QUEUE "meteostation:bmp280"

#define POST_FREQUENCY 500


#define BMP_SCK  (13)
#define BMP_MISO (12)
#define BMP_MOSI (11)
#define BMP_CS   (10)



Adafruit_BMP280 bmp; // use I2C interface
//Adafruit_BMP280 bmp(BMP_CS); // hardware SPI
//Adafruit_BMP280 bmp(BMP_CS, BMP_MOSI, BMP_MISO,  BMP_SCK);

Preferences preferences;

Adafruit_Sensor *bmp_temperature = bmp.getTemperatureSensor();
Adafruit_Sensor *bmp_pressure = bmp.getPressureSensor();


// Replace with your network credentials
const char* ssid = "BSD";
const char* password = "kuku01KUKU01";

WiFiClient redisConn;
Redis *gRedis = nullptr;
DNSClient dns;

void setup() {
  Serial.begin(115200);
  Serial.println(F("BMP280 Sensor event test"));
  
  //if (!bmp.begin(BMP280_ADDRESS_ALT, BMP280_CHIPID)) {
  if (!bmp.begin()) {
      Serial.println(F("Could not find a valid BMP280 sensor, check wiring or "
                        "try a different address!"));
      while (1) delay(10);
  }

  /* Default settings from datasheet. */
  bmp.setSampling(Adafruit_BMP280::MODE_NORMAL,     /* Operating Mode. */
                  Adafruit_BMP280::SAMPLING_X2,     /* Temp. oversampling */
                  Adafruit_BMP280::SAMPLING_X16,    /* Pressure oversampling */
                  Adafruit_BMP280::FILTER_X16,      /* Filtering. */
                  Adafruit_BMP280::STANDBY_MS_500); /* Standby time. */

  bmp_temp->printSensorDetails();


/*
  // get/set settings from/to flash storage. Reboot and poweroff are not clear these settings
  uint32_t redis_host;
  String ssid, password, redis_host, redis_pass;
  preferences.begin("settings", false);

  preferences.putString("ssid", ssid);
  preferences.putString("password", password);
  preferences.putString("redis_host", redis_host);
  preferences.putUInt("redis_port", redis_port);
  preferences.putString("redis_pass", redis_pass);
  preferences.end();


  preferences.begin("settings", false);
  ssid = preferences.getString("ssid", ""); 
  password = preferences.getString("password", "");
  redis_host = preferences.getString("redis_host", "");
  redis_port = preferences.getUInt("redis_port", "");
  redis_pass = preferences.getString("redis_pass", "");
  preferences.end();
*/


  // Connect to Wi-Fi network with SSID and password
  Serial.print("Connecting to ");
  Serial.println(ssid);
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
      delay(500);
      Serial.print(".");
  }

  // Print local IP address and start web server
  Serial.println("");
  Serial.println("WiFi connected.");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());

  if (!redisConn.connect(REDIS_ADDR, REDIS_PORT))
  {
      Serial.println("Failed to connect to the Redis server!");
      return;
  }

  gRedis = new Redis(redisConn);
  auto connRet = gRedis->authenticate(REDIS_PASSWORD);
  if (connRet == RedisSuccess)
  {
      Serial.printf("Connected to the Redis server at %s!\n", REDIS_ADDR);
  }
  else
  {
      Serial.printf("Failed to authenticate to the Redis server! Errno: %d\n", (int)connRet);
      return;
  }

  //  keep data in queue specified time in seconds. Default: one month statistics
  //gRedis->expire(REDIS_QUEUE, 31 * 24 * 60 * 60);

}

StaticJsonDocument<2048> doc;
unsigned long lastPost = 0;

void loop() {

  auto startTime = millis();
  if (startTime - lastPost > POST_FREQUENCY)
  {

/*
// Sensor event (36 bytes)
// struct sensor_event_s is used to provide a single sensor event in a common format.
typedef struct {
  int32_t version;   //< must be sizeof(struct sensors_event_t)
  int32_t sensor_id; // unique sensor identifier
  int32_t type;      // sensor type
  int32_t reserved0; // reserved
  int32_t timestamp; // time is in milliseconds
  union {
    float data[4];              ///< Raw data
    sensors_vec_t acceleration; // acceleration values are in meter per second
                                   per second (m/s^2)
    sensors_vec_t
        magnetic; // magnetic vector values are in micro-Tesla (uT)
    sensors_vec_t orientation; // orientation values are in degrees
    sensors_vec_t gyro;        // gyroscope values are in rad/s
    float temperature; // temperature is in degrees centigrade (Celsius)
    float distance;    // distance in centimeters
    float light;       // light in SI lux units
    float pressure;    // pressure in hectopascal (hPa)
    float relative_humidity; // relative humidity in percent
    float current;           // current in milliamps (mA)
    float voltage;           // voltage in volts (V)
    sensors_color_t color;   // color in RGB component values
  };                         ///< Union for the wide ranges of data we can carry
} sensors_event_t;
*/

    sensors_event_t temp_event, pressure_event;
    bmp_temperature->getEvent(&temp_event);
    bmp_pressure->getEvent(&pressure_event);

    Serial.print(F("Temperature = "));
    Serial.print(temp_event.temperature);
    Serial.println(" *C");
  
    Serial.print(F("Pressure = "));
    Serial.print(pressure_event.pressure);
    Serial.println(" hPa");
    Serial.println();


    // push data to redis
    doc["time"] = temp_event.timestamp;
    doc["temperature"] = temp_event.temperature;
    doc["pressure"] = pressure_event.pressure;

    String jsonStr;
    serializeJson(doc, jsonStr);
    Serial.printf("Sending JSON payload:\n\t'%s'\n", jsonStr.c_str());

    //auto listeners = gRedis->publish("arduino-redis:jsonpub", jsonStr.c_str());
    //auto listeners = gRedis->publish("meteostation:bmp280", jsonStr.c_str());
    auto list_length = gRedis->lpush(REDIS_QUEUE, jsonStr.c_str());

    
    Serial.printf("Push sensor values to %s list in redis storage: %d\n", REDIS_QUEUE, list_length);

    doc.clear();
    lastPost = millis();
  }

}
