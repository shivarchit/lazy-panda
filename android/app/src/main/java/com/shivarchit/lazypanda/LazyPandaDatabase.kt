package com.shivarchit.lazypanda

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase

@Database(entities = [TokenModel::class], version = 1)
abstract class LazyPandaDatabase : RoomDatabase() {
    abstract fun tokenModel(): TokenModel

    companion object {
        private var INSTANCE: LazyPandaDatabase? = null

        fun getDatabaseInstance(context: Context): LazyPandaDatabase {
            return INSTANCE ?: synchronized(this) {
                val instance = Room.databaseBuilder(
                    context.applicationContext,
                    LazyPandaDatabase::class.java,
                    "lazy_panda"
                ).build()
                INSTANCE = instance
                instance
            }
        }
    }
}