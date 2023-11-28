package com.shivarchit.lazypanda

import android.content.Intent
import android.os.Bundle
import android.widget.Button
import android.widget.EditText
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import retrofit2.Retrofit
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import retrofit2.converter.gson.GsonConverterFactory
import android.app.Application

class MyApp : Application() {
    companion object {
        lateinit var apiService: ApiService
        lateinit var token: String
    }

    override fun onCreate() {
        super.onCreate()

        val retrofit = Retrofit.Builder()
            .baseUrl("https://starfish-hopeful-spaniel.ngrok-free.app/")
            .addConverterFactory(GsonConverterFactory.create())
            .build()

        apiService = retrofit.create(ApiService::class.java)
    }
}

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        val loginButton = findViewById<Button>(R.id.loginButton)
        loginButton.setOnClickListener {
            val usernameEditText = findViewById<EditText>(R.id.username)
            val passwordEditText = findViewById<EditText>(R.id.password)

            val username = usernameEditText.text.toString()
            val password = passwordEditText.text.toString()

            if (username.isEmpty() || password.isEmpty()) {
                showToast("Please enter both username and password.")
            } else {
                loginUser(username, password)
            }
        }
    }

    private fun navigateToTrackpadView() {
        val intent = Intent(this, TrackpadActivity::class.java)
        startActivity(intent)
        finish()
    }

    private fun showToast(message: String) {
        Toast.makeText(this, message, Toast.LENGTH_SHORT).show()
    }

    private fun loginUser(username: String, password: String) {
        val call = MyApp.apiService.login(username, password)
        call.enqueue(object : Callback<ResponseModel> {
            override fun onResponse(call: Call<ResponseModel>, response: Response<ResponseModel>) {
                if (response.isSuccessful) {
                    val responseBody = response.body()
                    if (responseBody != null && responseBody.success) {
                        responseBody.token?.let {
                            MyApp.token = it
                            navigateToTrackpadView()
                        } ?: showToast("Token not found in the response.")
                    } else {
                        showToast(responseBody?.message ?: "Login failed.")
                    }
                } else {
                    showToast("Failed to make the login request.")
                }
            }

            override fun onFailure(call: Call<ResponseModel>, t: Throwable) {
                showToast("Network error: ${t.message}")
            }
        })
    }
}
