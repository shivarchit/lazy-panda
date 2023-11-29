package com.shivarchit.lazypanda

// ApiService.kt
import retrofit2.Call
import retrofit2.http.Field
import retrofit2.http.Headers
import retrofit2.http.POST
import retrofit2.http.Body

interface ApiService {
    @POST("api/login")
    fun login(@Body body: Map<String, String>
    ): Call<ResponseModel>
}
