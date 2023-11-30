package com.shivarchit.lazypanda

import UserTokenDao
import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase

@Database(entities = [UserToken::class], version = 1)
abstract class LazyPandaDatabase : RoomDatabase() {
    abstract fun userTokenDao(): UserTokenDao
}