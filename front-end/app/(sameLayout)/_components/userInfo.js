import { timeAgo } from '@/app/helpers'
import React from 'react'

const UserInfo = ({ userInfo, articleInfo }) => {
    return (
        <div className="user">
            <div className="profile-photo">
                <img src="./images/profile-13.jpg" />
            </div>
            <div className="ingo">
                <h3>{userInfo.first_name} {userInfo.last_name}</h3>
                {articleInfo && articleInfo.parent ==null &&
                    <>
                        <small>{articleInfo.privacy}</small> . <small>{timeAgo(articleInfo.created_at)}</small>
                    </>
                }
            </div>
        </div>
    )
}

export default UserInfo
