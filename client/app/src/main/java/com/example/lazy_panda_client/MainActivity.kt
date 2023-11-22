package com.shivarchit.lazy_panda_client

import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.TimeoutCancellationException
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import kotlinx.coroutines.withTimeout
//import kotlinx.serialization.Serializable
//import kotlinx.serialization.decodeFromString
//import kotlinx.serialization.json.Json
import org.json.JSONObject
import org.w3c.dom.Text
import java.io.BufferedReader
import java.io.BufferedWriter
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

class MainActivity : AppCompatActivity() {

    private lateinit var button: Button;
    private lateinit var textView: TextView;

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        button = findViewById<Button>(R.id.space) as Button
        textView = findViewById<TextView>(R.id.textView) as TextView
        button.setOnClickListener {
            Toast.makeText(
                this@MainActivity,
                "${button.text} key pressed, sending request",
                Toast.LENGTH_SHORT
            ).show()
            lifecycleScope.launch {
                httpRequesterOnClickHandler(button.text.toString())
            }
        }
    }

    private suspend fun httpRequesterOnClickHandler(buttonText: String) {
        val jsonData = JSONObject()
        jsonData.put("key", buttonText.toString().lowercase())

        val url = URL("http://10.0.2.2:3010/keyboard-event")

        try {
            // Use withTimeout to set a timeout for the HTTP request
            val response = withTimeout(10000) {
                withContext(Dispatchers.IO) {
                    val urlConnection = url.openConnection() as HttpURLConnection
                    urlConnection.requestMethod = "POST"
                    urlConnection.setRequestProperty("Content-Type", "application/json; utf-8")
                    urlConnection.doOutput = true

                    val outputStream = urlConnection.outputStream
                    val writer = BufferedWriter(OutputStreamWriter(outputStream, "UTF-8"))
                    writer.write(jsonData.toString())
                    writer.flush()
                    writer.close()
                    outputStream.close()

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
            }

            withContext(Dispatchers.Main) {
                // Update the UI with a Toast message indicating success and the response
                Toast.makeText(this@MainActivity, "Successful $response", Toast.LENGTH_LONG).show()
            }
        } catch (e: TimeoutCancellationException) {
            // Handle timeout exception
            withContext(Dispatchers.Main) {
                println(e)
                Toast.makeText(this@MainActivity, "Request timed out", Toast.LENGTH_SHORT).show()
                textView.text = "Error $e"
            }
        } catch (e: Exception) {
            // Handle other exceptions
            withContext(Dispatchers.Main) {
                println(e)
                Toast.makeText(this@MainActivity, "Error: ${e.message}", Toast.LENGTH_SHORT).show()
                textView.text = "Error $e"
            }
        }
    }


}
