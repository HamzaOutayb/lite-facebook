// "use client";
// import { createContext, useContext, useEffect, useRef, useState } from "react";

// const WorkerContext = createContext(
//     {
//         portRef: { current: null }
//     }
// );

// export function WorkerProvider({ children }) {
//     const workerRef = useRef(null);
//     const portRef = useRef(null);
//     const [clientWorker, setClientWorker] = useState(null);
//     const [conversations, setConversations] = useState([]);
//     const selectedConversationRef = useRef(null);
//     const [messages, setMessages] = useState([]);
//     const userRef = useRef(null);
//     const [notfications, setNotifications] = useState(0);

//     useEffect(() => {
//         const storedUser = localStorage.getItem("user");
//         console.log("storedUser", storedUser);
//         userRef.current = storedUser ? JSON.parse(storedUser) : null;
//         console.log("userRef", userRef.current);
//     }, []);

//     useEffect(() => {
//         const worker = new SharedWorker("/sharedworker.js");
//         workerRef.current = worker;
//         portRef.current = worker.port;
//         portRef.current.start();
//         setClientWorker(worker);

//         const port = portRef.current;
//         if (!port) return;
//         const messageHandler = ({ data }) => {
//             switch (data.type) {
//                 case "conversations":
//                     const onlineUsers = data.online_users;
//                     if (data.conversations && data.conversations.length > 0) {
//                         setConversations(
//                             data.conversations?.map((conv) => {
//                                 if (conv.user_info) {
//                                     return {
//                                         ...conv,
//                                         user_info: {
//                                             ...conv.user_info,
//                                             online: onlineUsers?.includes(conv.user_info.id),
//                                         },
//                                     };
//                                 }
//                                 return conv;
//                             })
//                         );
//                     }

//                     break;

//                 case "online":
//                 case "offline":
//                     setConversations((prev) =>
//                         prev.map((conv) => {
//                             if (conv.user_info?.id === data.user_info.id) {
//                                 return {
//                                     ...conv,
//                                     user_info: {
//                                         ...conv.user_info,
//                                         online: data.type === "online",
//                                     },
//                                 };
//                             }
//                             return conv;
//                         })
//                     );
//                     break;

//                 case "new_message":
//                     const msg = data.message;
//                     const conversationId = msg.conversation_id;

//                     setConversations((prev) => {
//                         const conversation = prev.find((c) => c.conversation.id === conversationId);
//                         if (conversation) {
//                             return [{
//                                 ...conversation,
//                                 last_message: data?.message?.content,
//                                 seen: conversation.conversation.id === selectedConversationRef.current?.id ? 0 : conversation.seen + 1,
//                             }, ...prev.filter((c) => c.conversation.id !== conversationId)];
//                         } else {
//                             return [
//                                 {
//                                     conversation: { id: msg.conversation_id },
//                                     user_info: { ...data.user_info, online: true },
//                                     last_message: data?.message?.content,
//                                     seen: 1
//                                 },
//                                 ...prev,
//                             ];
//                         }
//                     });

//                     console.log("selectedConversationRef", selectedConversationRef.current);
//                     console.log("conversationId", conversationId);
//                     if (selectedConversationRef.current?.id === conversationId) {
//                         setMessages((prev) => [...prev, data]);
//                         console.log("userRef.current => ", userRef.current);
//                         if (userRef.current && userRef.current.id !== data.user_info.id) {
//                             // console.log("user", user);
//                             // console.log("data", data);
//                             const type = selectedConversationRef.current.type == "private" ? "read_messages_private" : "read_messages_group";
//                             port.postMessage({
//                                 kind: "send",
//                                 payload: {
//                                     type,
//                                     message: {
//                                         conversation_id: conversationId,
//                                     },
//                                 },
//                             });
//                         }
//                     }
//                     break;

//                 default:
//                     console.warn("Unhandled message type:", data.type);
//             }
//         };


//         port.addEventListener("message", messageHandler);
//         port.postMessage({ kind: "connect", payload: process.env.NEXT_PUBLIC_API_URL });

//         return () => {
//             port.removeEventListener("message", messageHandler);
//             port?.postMessage({
//                 kind: "close"
//             })
//             portRef.current?.close();
//             portRef.current = null;
//             workerRef.current = null;
//         };
//     }, []);

//     useEffect(() => {
//         const count = conversations?.reduce((acc, conv) => acc + conv.seen, 0);
//         setNotifications(count);
//         console.log(count)
//     }, [conversations])

//     const value = {
//         portRef,
//         userRef,
//         clientWorker,
//         conversations,
//         setConversations,
//         selectedConversationRef,
//         messages,
//         setMessages,
//         notfications
//     }

//     return (
//         <WorkerContext.Provider
//             value={value}
//         >
//             {children}
//         </WorkerContext.Provider>
//     );
// }

// export const useWorker = () => useContext(WorkerContext);










////////////////////

"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import FakePort from "./FakePort";

const WorkerContext = createContext({
    portRef: { current: null }
});

export function WorkerProvider({ children }) {
    const workerRef = useRef(null);
    const portRef = useRef(null);
    const [clientWorker, setClientWorker] = useState(null);
    const [conversations, setConversations] = useState([]);
    const selectedConversationRef = useRef(null);
    const [messages, setMessages] = useState([]);
    const userRef = useRef(null);
    const [notifications, setNotifications] = useState(0);

    // Load user from localStorage
    useEffect(() => {
        const storedUser = localStorage.getItem("user");
        console.log("storedUser", storedUser);
        userRef.current = storedUser ? JSON.parse(storedUser) : null;
        console.log("userRef", userRef.current);
    }, []);

    // Initialize worker or fallback
    useEffect(() => {
        let port;

        // Check if SharedWorker is supported
        if (typeof SharedWorker !== "undefined") {
            const worker = new SharedWorker("/sharedworker.js");
            workerRef.current = worker;
            port = worker.port;
            port.start();
            setClientWorker(worker); // Optional: keep for compatibility if used elsewhere
        } else {
            // Fallback to FakePort with WebSocket
            port = new FakePort(process.env.NEXT_PUBLIC_API_URL);
            setClientWorker(null); // No worker in this case
        }

        portRef.current = port;

        const messageHandler = ({ data }) => {
            switch (data.type) {
                case "conversations":
                    const onlineUsers = data.online_users;
                    if (data.conversations && data.conversations.length > 0) {
                        setConversations(
                            data.conversations.map((conv) => {
                                if (conv.user_info) {
                                    return {
                                        ...conv,
                                        user_info: {
                                            ...conv.user_info,
                                            online: onlineUsers?.includes(conv.user_info.id),
                                        },
                                    };
                                }
                                return conv;
                            })
                        );
                    }
                    break;

                case "online":
                case "offline":
                    setConversations((prev) =>
                        prev.map((conv) => {
                            if (conv.user_info?.id === data.user_info.id) {
                                return {
                                    ...conv,
                                    user_info: {
                                        ...conv.user_info,
                                        online: data.type === "online",
                                    },
                                };
                            }
                            return conv;
                        })
                    );
                    break;

                case "new_message":
                    const msg = data.message;
                    const conversationId = msg.conversation_id;

                    setConversations((prev) => {
                        const conversation = prev.find((c) => c.conversation.id === conversationId);
                        if (conversation) {
                            return [{
                                ...conversation,
                                last_message: data?.message?.content,
                                seen: userRef?.current?.id != msg.sender_id ? (selectedConversationRef.current?.id === conversationId ? conversation.seen : conversation.seen + 1) : 0,
                                // seen: conversation.conversation.id === selectedConversationRef.current?.id ? 0 : conversation.seen + 1,
                            }, ...prev.filter((c) => c.conversation.id !== conversationId)];
                        } else {
                            return [
                                {
                                    conversation: { id: msg.conversation_id },
                                    user_info: { ...data.user_info, online: true },
                                    last_message: data?.message?.content,
                                    seen: userRef?.current?.id != msg.sender_id ? 1 : 0,
                                    // senn : 1
                                },
                                ...prev,
                            ];
                        }
                    });

                    if (selectedConversationRef.current?.id === conversationId) {
                        setMessages((prev) => [...prev, data]);
                        if (userRef.current && userRef.current.id !== data.user_info.id) {
                            const type = selectedConversationRef.current.type === "private" ? "read_messages_private" : "read_messages_group";
                            port.postMessage({
                                kind: "send",
                                payload: {
                                    type,
                                    message: { conversation_id: conversationId },
                                },
                            });
                        }
                    }
                    break;

                default:
                    console.warn("Unhandled message type:", data.type);
            }
        };

        port.addEventListener("message", messageHandler);
        port.postMessage({ kind: "connect", payload: process.env.NEXT_PUBLIC_API_URL });

        return () => {
            port.removeEventListener("message", messageHandler);
            port?.postMessage({
                kind: "close"
            })
            if (port.close) {
                port.close();
            }
            portRef.current = null;
            workerRef.current = null;
        };
    }, []);

    useEffect(() => {
        const count = conversations.reduce((acc, conv) => acc + conv.seen, 0);
        setNotifications(count);
        console.log(count);
    }, [conversations]);

    const value = {
        portRef,
        userRef,
        clientWorker,
        conversations,
        setConversations,
        selectedConversationRef,
        messages,
        setMessages,
        notifications
    };

    return (
        <WorkerContext.Provider value={value}>
            {children}
        </WorkerContext.Provider>
    );
}

export const useWorker = () => useContext(WorkerContext);