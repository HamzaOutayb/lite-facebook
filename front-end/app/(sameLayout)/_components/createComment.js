import React, { useState } from 'react'
import styles from "./comment.module.css"
// import styles from "./createPostModal.module.css"
import { AddPhotoAlternate } from '@mui/icons-material'
import { addArticle } from '@/app/helpers'

const CreateComment = ({ setComments, setCommentCount, parent }) => {
    const [imagePreview, setImagePreview] = useState("")
    const [commentContent, setCommentContent] = useState("")


    return (
        <>
            <form
                className={styles.form}
                onSubmit={async (e) => {
                    const added = addArticle(e, setComments, { parent })
                    if (added) {
                        setCommentContent("")
                        setCommentCount((prev) => prev + 1)
                    }
                }}
            >
                <textarea
                    className={styles.textInput}
                    value={commentContent}
                    name='content'
                    onChange={(e) => setCommentContent(e.target.value)}
                    placeholder="Write a comment..."
                ></textarea>
                <input
                    type="file"
                    id='postImage'
                    name='image'
                    onChange={(e) => {
                        if (e.target.files[0]) {
                            const file = e.target.files[0]
                            const reader = new FileReader()
                            reader.onloadend = () => {
                                setImagePreview(reader.result)
                            }
                            reader.readAsDataURL(file)
                        }else{
                            setImagePreview("")
                        }
                    }}
                    className={styles.fileInput} />
                <div className={styles.actionButtons}>
                    <label htmlFor="postImage" className={styles.addFile}>

                        {imagePreview ? <img src={imagePreview} className="imagePreview" /> : <AddPhotoAlternate />}

                    </label>
                    <button type="submit" className="btn btn-primary">Comment</button>
                </div>
            </form>
        </>
    )
}

export default CreateComment
