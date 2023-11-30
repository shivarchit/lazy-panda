import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import androidx.room.Update
import com.shivarchit.lazypanda.UserToken

@Dao
interface UserTokenDao {
    @Query("SELECT * FROM user_token")
    fun getAllUserTokens(): List<UserToken>

    @Query("SELECT * FROM user_token WHERE id = :id")
    fun getUserTokenById(id: String): UserToken?

    @Insert
    fun insertUserToken(userToken: UserToken)

    @Update
    fun updateUserToken(userToken: UserToken)

    @Delete
    fun deleteUserToken(userToken: UserToken)
}