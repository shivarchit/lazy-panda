package com.shivarchit.lazy_panda_client

import android.os.Bundle
import android.view.MotionEvent
import android.view.View
import android.widget.Button
import android.widget.TextView
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
//import com.example.lazy_panda_client.R
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.TimeoutCancellationException
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import kotlinx.coroutines.withTimeout
//import kotlinx.serialization.Serializable
//import kotlinx.serialization.decodeFromString
//import kotlinx.serialization.json.Json
import org.json.JSONObject
import java.io.BufferedReader
import java.io.BufferedWriter
import java.io.IOException
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

class MainActivity : AppCompatActivity() {

    private lateinit var button: Button;
    private lateinit var textView: TextView;
    private  lateinit var trackPad: View;

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        button = findViewById<Button>(R.id.space) as Button
        textView = findViewById<TextView>(R.id.textView) as TextView
        trackPad = findViewById<View>(R.id.trackpad) as View

        trackPad.setOnTouchListener { _, event ->
            handleTouch(event)
        }
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

    private fun handleTouch(event: MotionEvent): Boolean {
        when (event.action) {
            MotionEvent.ACTION_MOVE -> {
                // Get the X and Y coordinates of the touch event
                val x = event.x
                val y = event.y

                // Update the TextView with the touch coordinates
                textView.text = "X: $x, Y: $y"
                Toast.makeText(this@MainActivity, "X: $x, Y: $y", Toast.LENGTH_SHORT).show()
                // Add your logic to simulate mouse movements based on touch coordinates
                // For example, you can send these coordinates to a server or perform other actions

                return true
            }
        }
        return false
    }

    private suspend fun httpRequesterOnClickHandler(buttonText: String) {
        val jsonData = JSONObject()
        jsonData.put("key", buttonText.toString().lowercase())

        val url = URL("https://starfish-hopeful-spaniel.ngrok-free.app/api/keyboard-event")

        try {
            val response = withTimeout(10000) {
                withContext(Dispatchers.IO) {
                    val urlConnection = url.openConnection() as HttpURLConnection
                    urlConnection.requestMethod = "POST"
                    urlConnection.setRequestProperty("Authorization","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE0MjQxMjUsInVzZXJfaWQiOiJzaGl2In0.VUqfM2wMuGN3cksQeTr2-aclk4jg9JBYg1op_I4jiIg")
                    urlConnection.setRequestProperty("Content-Type", "application/json; utf-8")
                    urlConnection.doOutput = true

                    val outputStream = urlConnection.outputStream
                    outputStream.use {
                        val writer = BufferedWriter(OutputStreamWriter(it, "UTF-8"))
                        writer.write(jsonData.toString())
                        writer.flush()
                    }

                    val httpResponseCode = urlConnection.responseCode
                    if (httpResponseCode == HttpURLConnection.HTTP_OK) {
                        val reader = BufferedReader(InputStreamReader(urlConnection.inputStream))
                        reader.use {
                            it.readText()
                        }
                    } else {
                        throw IOException("HTTP error code: $httpResponseCode")
                    }
                }
            }

            withContext(Dispatchers.Main) {
                Toast.makeText(this@MainActivity, "Successful $response", Toast.LENGTH_LONG).show()
            }
        } catch (e: TimeoutCancellationException) {
            withContext(Dispatchers.Main) {
                Toast.makeText(this@MainActivity, "Request timed out", Toast.LENGTH_SHORT).show()
                textView.text = "Error: Request timed out"
            }
        } catch (e: Exception) {
            withContext(Dispatchers.Main) {
                Toast.makeText(this@MainActivity, "Error: ${e.message}", Toast.LENGTH_SHORT).show()
                textView.text = "Error: ${e.message}"
            }
        }
    }



}
