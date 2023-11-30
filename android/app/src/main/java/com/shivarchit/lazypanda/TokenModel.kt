package com.shivarchit.lazypanda

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import java.util.UUID


@Entity(tableName = "user_token")
data class UserToken(
    @PrimaryKey val id: String = UUID.randomUUID().toString(),
    val username: String,
    val token: String,
    val expiryTime: Long
)
