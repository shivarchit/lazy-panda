package com.shivarchit.lazypanda

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import androidx.room.Update
import java.sql.Time

@Dao
interface DatabaseInterface {
//    @Query("SELECT * FROM user")
//    fun getAllTokens(): List<TokenModel>

//    @Query("SELECT * FROM token_v1 WHERE username LIKE :username")
//    fun findByUsername(username: String): TokenModel

    @Update
    fun updateUserToken(username: String,token: String,tokenExpiry: Time):TokenModel

    @Insert
    fun insertUserToken(vararg users: TokenModel)

    @Delete
    fun deleteUserToken(user: TokenModel)
}