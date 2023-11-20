package com.example.lazy_panda_client

import android.os.Bundle
import android.widget.Button
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.HttpURLConnection
import java.net.URL

class MainActivity : AppCompatActivity() {

    private lateinit var button: Button

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
        // Replace the URL with your actual endpoint
        val url = URL("http://localhost:3010/keyboard-event")

        try {
            // Use withContext to switch to the IO dispatcher for network operations
            val response = withContext(Dispatchers.IO) {
                val urlConnection = url.openConnection() as HttpURLConnection
                urlConnection.requestMethod = "GET"

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

            // Process the response on the UI thread
            withContext(Dispatchers.Main) {
                // Update UI with the response and button text
                val updatedText = "$buttonText Key Pressed, sending request\nResponse: $response"
                button.text = updatedText
            }

        } catch (e: Exception) {
            // Handle the exception on the UI thread
            withContext(Dispatchers.Main) {
                // Show error message or handle accordingly
                println(e)
                Toast.makeText(this@MainActivity, "Error: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }
}
