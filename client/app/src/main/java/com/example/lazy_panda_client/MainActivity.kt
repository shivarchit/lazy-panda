package com.example.lazy_panda_client

import android.os.Bundle
import android.widget.Button
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import org.json.JSONObject
import java.io.BufferedReader
import java.io.BufferedWriter
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

class MainActivity : AppCompatActivity() {

    private lateinit var button: Button;

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        button = findViewById<Button>(R.id.space) as Button
        button.setOnClickListener {
            Toast.makeText(this@MainActivity, "${button.text} key pressed, sending request", Toast.LENGTH_SHORT).show()
            lifecycleScope.launch {
                httpRequesterOnClickHandler(button.text.toString())
            }
        }
    }

    private suspend fun httpRequesterOnClickHandler(buttonText: String) {
        // Create a sample JSON data
        val jsonData = JSONObject()
        jsonData.put("key", buttonText.toString().lowercase())

        val url = URL("http://10.0.2.2:3010/keyboard-event")

        try {
            val response = withContext(Dispatchers.IO) {
                val urlConnection = url.openConnection() as HttpURLConnection
                urlConnection.requestMethod = "POST"

                // Set content type to application/json
                urlConnection.setRequestProperty("Content-Type", "application/json; utf-8")

                // Enable output stream for writing data
                urlConnection.doOutput = true

                // Write the JSON data to the output stream
                val outputStream = urlConnection.outputStream
                val writer = BufferedWriter(OutputStreamWriter(outputStream, "UTF-8"))
                writer.write(jsonData.toString())
                writer.flush()
                writer.close()
                outputStream.close()

                // Get the response
                val reader = BufferedReader(InputStreamReader(urlConnection.inputStream))
                val response = StringBuilder()

                reader.use {
                    var line = it.readLine()
                    while (line != null) {
                        response.append(line)
                        line = it.readLine()
                    }
                }

                response.toString()
            }

            withContext(Dispatchers.Main) {
                // Update UI with the response and button text
                Toast.makeText(this@MainActivity, "Successful", Toast.LENGTH_LONG).show()
            }

        } catch (e: Exception) {
            withContext(Dispatchers.Main) {
                println(e)
                Toast.makeText(this@MainActivity, "Error: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }


}
