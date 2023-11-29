package com.shivarchit.lazypanda

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import java.sql.Time
import java.util.UUID

@Entity
data class TokenModel(
    @PrimaryKey val uuid: UUID,
    @ColumnInfo(name = "username") val username: String?,
    @ColumnInfo(name = "token") val token: String?,
    @ColumnInfo(name = "tokenExpiry") val tokenExpiry: Time?
)
