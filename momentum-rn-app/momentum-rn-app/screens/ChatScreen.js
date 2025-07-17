// import React, { useState, useEffect, useRef } from 'react';
// import {
//   View,
//   Text,
//   TextInput,
//   TouchableOpacity,
//   ScrollView,
//   StyleSheet,
//   SafeAreaView,
//   Alert,
//   ActivityIndicator,
//   KeyboardAvoidingView,
//   Platform,
// } from 'react-native';
// import AsyncStorage from '@react-native-async-storage/async-storage';

// export default function ChatScreen({ userInfo, navigation }) {
//   const [messages, setMessages] = useState([]);
//   const [newMessage, setNewMessage] = useState('');
//   const [ws, setWs] = useState(null);
//   const [connected, setConnected] = useState(false);
//   const [recipientId, setRecipientId] = useState('');
//   const [activeChat, setActiveChat] = useState(null);
//   const [loading, setLoading] = useState(false);
//   const scrollViewRef = useRef();

//   useEffect(() => {
//     connectWebSocket();
//     return () => {
//       if (ws) {
//         ws.close();
//       }
//     };
//   }, []);

//   const connectWebSocket = async () => {
//     try {
//       const token = await AsyncStorage.getItem('token');
//       if (!token) {
//         Alert.alert('Error', 'No authentication token found');
//         return;
//       }

//       const wsUrl = `ws://10.0.2.2:8083/ws`;
//       const websocket = new WebSocket(wsUrl, [], {
//         headers: {
//           Authorization: `Bearer ${token}`,
//         },
//       });

//       websocket.onopen = () => {
//         console.log('WebSocket connected');
//         setConnected(true);
//       };

//       websocket.onmessage = (event) => {
//         const data = JSON.parse(event.data);
//         handleWebSocketMessage(data);
//       };

//       websocket.onerror = (error) => {
//         console.log('WebSocket error:', error);
//         Alert.alert('Connection Error', 'Failed to connect to chat service');
//       };

//       websocket.onclose = () => {
//         console.log('WebSocket disconnected');
//         setConnected(false);
//       };

//       setWs(websocket);
//     } catch (error) {
//       console.log('Error connecting WebSocket:', error);
//       Alert.alert('Error', 'Failed to connect to chat service');
//     }
//   };

//   const handleWebSocketMessage = (data) => {
//     switch (data.type) {
//       case 'chat_history':
//         // Initial chat history - combine read and unread
//         const allMessages = [...data.read, ...data.unread].sort(
//           (a, b) => new Date(a.created_at) - new Date(b.created_at)
//         );
//         setMessages(allMessages);
//         break;

//       case 'chat_with_user':
//         // Chat with specific user
//         const chatMessages = [...data.read, ...data.unread].sort(
//           (a, b) => new Date(a.created_at) - new Date(b.created_at)
//         );
//         setMessages(chatMessages);
//         setActiveChat(data.with);
//         break;

//       case 'new_message':
//         // New incoming message
//         const newMsg = {
//           sender_id: data.from,
//           receiver_id: userInfo.id,
//           content: data.content,
//           created_at: new Date().toISOString(),
//           is_read: true,
//         };
//         setMessages(prev => [...prev, newMsg]);
//         break;

//       case 'message_read':
//         // Update messages sent *by me* to this user, mark them as read
//         setMessages((prevMessages) =>
//           prevMessages.map((msg) => {
//             if (
//               msg.sender_id === userInfo.id &&
//               msg.receiver_id === data.by &&
//               !msg.is_read
//             ) {
//               return { ...msg, is_read: true };
//             }
//             return msg;
//           })
//         );
//         break;
  

//       case 'error':
//         Alert.alert('Error', data.error);
//         break;

//       default:
//         console.log('Unknown message type:', data.type);
//     }
//   };

//   const sendMessage = () => {
//     if (!newMessage.trim() || !recipientId || !ws || !connected) {
//       Alert.alert('Error', 'Please enter a message and recipient ID');
//       return;
//     }

//     const messageData = {
//       type: 'message',
//       receiver_id: parseInt(recipientId),
//       content: newMessage.trim(),
//     };

//     ws.send(JSON.stringify(messageData));

//     // Add message to local state optimistically
//     const sentMsg = {
//       sender_id: userInfo.id,
//       receiver_id: parseInt(recipientId),
//       content: newMessage.trim(),
//       created_at: new Date().toISOString(),
//       is_read: false,
//     };
//     setMessages(prev => [...prev, sentMsg]);
//     setNewMessage('');
//   };

//   const loadChatWithUser = () => {
//     if (!recipientId || !ws || !connected) {
//       Alert.alert('Error', 'Please enter a valid user ID');
//       return;
//     }

//     setLoading(true);
//     const loadData = {
//       type: 'load_chat',
//       with: parseInt(recipientId),
//     };

//     ws.send(JSON.stringify(loadData));
//     setLoading(false);
//   };

//   const checkMutualFollow = async () => {
//     if (!recipientId) {
//       Alert.alert('Error', 'Please enter a user ID');
//       return;
//     }

//     try {
//       const token = await AsyncStorage.getItem('token');
//       const response = await fetch('http://10.0.2.2:8083/check-mutual-follow', {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//           Authorization: `Bearer ${token}`,
//         },
//         body: JSON.stringify({ other_user_id: parseInt(recipientId) }),
//       });

//       const data = await response.json();
//       if (response.ok) {
//         Alert.alert(
//           'Mutual Follow Status',
//           data.mutual ? 'You are mutual followers!' : 'You are not mutual followers'
//         );
//       } else {
//         Alert.alert('Error', 'Failed to check mutual follow status');
//       }
//     } catch (error) {
//       console.log('Error checking mutual follow:', error);
//       Alert.alert('Error', 'Failed to check mutual follow status');
//     }
//   };

//   const followUser = async () => {
//     if (!recipientId) {
//       Alert.alert('Error', 'Please enter a user ID');
//       return;
//     }

//     try {
//       const token = await AsyncStorage.getItem('token');
//       const response = await fetch('http://10.0.2.2:8083/follow', {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//           Authorization: `Bearer ${token}`,
//         },
//         body: JSON.stringify({ other_user_id: parseInt(recipientId) }),
//       });

//       const data = await response.json();
//       if (response.ok) {
//         Alert.alert('Success', 'User followed successfully');
//       } else {
//         Alert.alert('Error', 'Failed to follow user');
//       }
//     } catch (error) {
//       console.log('Error following user:', error);
//       Alert.alert('Error', 'Failed to follow user');
//     }
//   };

//   const formatTime = (timestamp) => {
//     const date = new Date(timestamp);
//     return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
//   };

//   const isMyMessage = (message) => {
//     return message.sender_id === userInfo.id;
//   };

//   return (
//     <SafeAreaView style={styles.container}>
//       <KeyboardAvoidingView 
//         style={styles.container} 
//         behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
//       >
//         {/* Header */}
//         <View style={styles.header}>
//           <TouchableOpacity 
//             style={styles.backButton}
//             onPress={() => navigation.goBack()}
//           >
//             <Text style={styles.backButtonText}>← Back</Text>
//           </TouchableOpacity>
//           <Text style={styles.title}>Chat</Text>
//           <View style={styles.connectionStatus}>
//             <View style={[styles.statusDot, { backgroundColor: connected ? '#4CAF50' : '#F44336' }]} />
//             <Text style={styles.statusText}>{connected ? 'Connected' : 'Disconnected'}</Text>
//           </View>
//         </View>

//         {/* User Controls */}
//         <View style={styles.controls}>
//           <View style={styles.inputRow}>
//             <TextInput
//               style={styles.userInput}
//               placeholder="Enter User ID to chat with"
//               value={recipientId}
//               onChangeText={setRecipientId}
//               keyboardType="numeric"
//             />
//             <TouchableOpacity style={styles.loadButton} onPress={loadChatWithUser}>
//               <Text style={styles.loadButtonText}>Load Chat</Text>
//             </TouchableOpacity>
//           </View>
          
//           <View style={styles.actionButtons}>
//             <TouchableOpacity style={styles.actionButton} onPress={checkMutualFollow}>
//               <Text style={styles.actionButtonText}>Check Mutual</Text>
//             </TouchableOpacity>
//             <TouchableOpacity style={styles.actionButton} onPress={followUser}>
//               <Text style={styles.actionButtonText}>Follow</Text>
//             </TouchableOpacity>
//           </View>
//         </View>

//         {/* Messages */}
//         <ScrollView
//           ref={scrollViewRef}
//           style={styles.messagesContainer}
//           onContentSizeChange={() => scrollViewRef.current?.scrollToEnd()}
//         >
//           {loading && (
//             <ActivityIndicator size="large" color="#007AFF" style={styles.loader} />
//           )}
          
//           {messages.map((message, index) => (
//             <View
//               key={index}
//               style={[
//                 styles.messageBubble,
//                 isMyMessage(message) ? styles.myMessage : styles.theirMessage,
//               ]}
//             >
//               <Text style={styles.messageText}>{message.content}</Text>
//               <Text style={styles.messageTime}>
//                 {formatTime(message.created_at)}
//                 {!message.is_read && isMyMessage(message) && (
//                   <Text style={styles.unreadIndicator}> • Unread</Text>
//                 )}
//               </Text>
//             </View>
//           ))}
//         </ScrollView>

//         {/* Message Input */}
//         <View style={styles.inputContainer}>
//           <TextInput
//             style={styles.messageInput}
//             placeholder="Type a message..."
//             value={newMessage}
//             onChangeText={setNewMessage}
//             multiline
//             maxLength={500}
//           />
//           <TouchableOpacity
//             style={[styles.sendButton, { opacity: newMessage.trim() ? 1 : 0.5 }]}
//             onPress={sendMessage}
//             disabled={!newMessage.trim()}
//           >
//             <Text style={styles.sendButtonText}>Send</Text>
//           </TouchableOpacity>
//         </View>
//       </KeyboardAvoidingView>
//     </SafeAreaView>
//   );
// }

// const styles = StyleSheet.create({
//   container: {
//     flex: 1,
//     backgroundColor: '#f5f5f5',
//   },
//   header: {
//     flexDirection: 'row',
//     justifyContent: 'space-between',
//     alignItems: 'center',
//     padding: 15,
//     backgroundColor: 'white',
//     borderBottomWidth: 1,
//     borderBottomColor: '#e0e0e0',
//   },
//   backButton: {
//     padding: 5,
//   },
//   backButtonText: {
//     fontSize: 16,
//     color: '#007AFF',
//   },
//   title: {
//     fontSize: 20,
//     fontWeight: 'bold',
//     color: '#333',
//   },
//   connectionStatus: {
//     flexDirection: 'row',
//     alignItems: 'center',
//   },
//   statusDot: {
//     width: 8,
//     height: 8,
//     borderRadius: 4,
//     marginRight: 5,
//   },
//   statusText: {
//     fontSize: 12,
//     color: '#666',
//   },
//   controls: {
//     backgroundColor: 'white',
//     padding: 15,
//     borderBottomWidth: 1,
//     borderBottomColor: '#e0e0e0',
//   },
//   inputRow: {
//     flexDirection: 'row',
//     marginBottom: 10,
//   },
//   userInput: {
//     flex: 1,
//     borderWidth: 1,
//     borderColor: '#ddd',
//     borderRadius: 8,
//     padding: 10,
//     marginRight: 10,
//     fontSize: 16,
//   },
//   loadButton: {
//     backgroundColor: '#007AFF',
//     paddingHorizontal: 15,
//     paddingVertical: 10,
//     borderRadius: 8,
//     justifyContent: 'center',
//   },
//   loadButtonText: {
//     color: 'white',
//     fontWeight: 'bold',
//   },
//   actionButtons: {
//     flexDirection: 'row',
//     justifyContent: 'space-around',
//   },
//   actionButton: {
//     backgroundColor: '#34C759',
//     paddingHorizontal: 20,
//     paddingVertical: 8,
//     borderRadius: 8,
//   },
//   actionButtonText: {
//     color: 'white',
//     fontWeight: 'bold',
//   },
//   messagesContainer: {
//     flex: 1,
//     padding: 15,
//   },
//   loader: {
//     marginVertical: 20,
//   },
//   messageBubble: {
//     maxWidth: '80%',
//     padding: 12,
//     borderRadius: 18,
//     marginBottom: 10,
//   },
//   myMessage: {
//     alignSelf: 'flex-end',
//     backgroundColor: '#007AFF',
//   },
//   theirMessage: {
//     alignSelf: 'flex-start',
//     backgroundColor: '#E5E5EA',
//   },
//   messageText: {
//     fontSize: 16,
//     color: 'white',
//   },
//   messageTime: {
//     fontSize: 12,
//     color: 'rgba(255, 255, 255, 0.7)',
//     marginTop: 4,
//   },
//   unreadIndicator: {
//     color: '#FFD60A',
//   },
//   inputContainer: {
//     flexDirection: 'row',
//     padding: 15,
//     backgroundColor: 'white',
//     alignItems: 'flex-end',
//   },
//   messageInput: {
//     flex: 1,
//     borderWidth: 1,
//     borderColor: '#ddd',
//     borderRadius: 20,
//     paddingHorizontal: 15,
//     paddingVertical: 10,
//     marginRight: 10,
//     fontSize: 16,
//     maxHeight: 100,
//   },
//   sendButton: {
//     backgroundColor: '#007AFF',
//     paddingHorizontal: 20,
//     paddingVertical: 10,
//     borderRadius: 20,
//   },
//   sendButtonText: {
//     color: 'white',
//     fontWeight: 'bold',
//   },
// });

// screens/ChatScreen.js
import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  StyleSheet,
  SafeAreaView,
  Alert,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  FlatList,
  RefreshControl,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';

export default function ChatScreen({ userInfo, navigation }) {
  // State management
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [ws, setWs] = useState(null);
  const [connected, setConnected] = useState(false);
  const [activeChat, setActiveChat] = useState(null);
  const [loading, setLoading] = useState(false);
  const [currentView, setCurrentView] = useState('chatList'); // 'chatList', 'userSearch', 'chat'
  const [refreshing, setRefreshing] = useState(false);
  
  // Chat list and users
  const [chatList, setChatList] = useState([]);
  const [allUsers, setAllUsers] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [followStatuses, setFollowStatuses] = useState({});
  
  const scrollViewRef = useRef();

  useEffect(() => {
    connectWebSocket();
    loadChatList();
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, []);

  const connectWebSocket = async () => {
    try {
      const token = await AsyncStorage.getItem('token');
      if (!token) {
        Alert.alert('Error', 'No authentication token found');
        return;
      }

      const wsUrl = `ws://aarya.live:8083/ws`;
      const websocket = new WebSocket(wsUrl, [], {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      websocket.onopen = () => {
        console.log('WebSocket connected');
        setConnected(true);
      };

      websocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleWebSocketMessage(data);
      };

      websocket.onerror = (error) => {
        console.log('WebSocket error:', error);
        Alert.alert('Connection Error', 'Failed to connect to chat service');
      };

      websocket.onclose = () => {
        console.log('WebSocket disconnected');
        setConnected(false);
      };

      setWs(websocket);
    } catch (error) {
      console.log('Error connecting WebSocket:', error);
      Alert.alert('Error', 'Failed to connect to chat service');
    }
  };

  const handleWebSocketMessage = (data) => {
    switch (data.type) {
      case 'connected':
        console.log('WebSocket connection confirmed for user:', data.user_id);
        break;

      case 'chat_with_user':
        // Chat with specific user
        const chatMessages = [...data.read, ...data.unread].sort(
          (a, b) => new Date(a.created_at) - new Date(b.created_at)
        );
        setMessages(chatMessages);
        setActiveChat({
          id: data.with,
          name: data.user_name,
          email: data.user_email,
        });
        setCurrentView('chat');
        break;

      case 'new_message':
        // New incoming message
        const newMsg = {
          sender_id: data.from,
          receiver_id: userInfo.id,
          content: data.content,
          created_at: new Date().toISOString(),
          is_read: true,
        };
        setMessages(prev => [...prev, newMsg]);
        // Refresh chat list to update unread counts
        loadChatList();
        break;

      case 'message_read':
        // Update messages sent by me to this user, mark them as read
        setMessages((prevMessages) =>
          prevMessages.map((msg) => {
            if (
              msg.sender_id === userInfo.id &&
              msg.receiver_id === data.by &&
              !msg.is_read
            ) {
              return { ...msg, is_read: true };
            }
            return msg;
          })
        );
        break;

      case 'message_sent':
        // Message sent confirmation
        console.log('Message sent successfully to user:', data.to);
        break;

      case 'error':
        Alert.alert('Error', data.error);
        break;

      default:
        console.log('Unknown message type:', data.type);
    }
  };

  const loadChatList = async () => {
    try {
      const token = await AsyncStorage.getItem('token');
      const response = await fetch('http://aarya.live:8083/chat-list', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        setChatList(data || []);
      }
    } catch (error) {
      console.log('Error loading chat list:', error);
    }
  };

  const loadAllUsers = async () => {
    try {
      const token = await AsyncStorage.getItem('token');
      const response = await fetch('http://aarya.live:8083/users', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        setAllUsers(data || []);
        // Load follow statuses for all users
        data.forEach(user => loadFollowStatus(user.id));
      }
    } catch (error) {
      console.log('Error loading users:', error);
    }
  };

  const searchUsers = async (query) => {
    if (!query.trim()) {
      setSearchResults([]);
      return;
    }

    try {
      const token = await AsyncStorage.getItem('token');
      const response = await fetch(`http://aarya.live:8083/users/search?q=${encodeURIComponent(query)}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        setSearchResults(data || []);
        // Load follow statuses for search results
        data.forEach(user => loadFollowStatus(user.id));
      }
    } catch (error) {
      console.log('Error searching users:', error);
    }
  };

  const loadFollowStatus = async (userId) => {
    try {
      const token = await AsyncStorage.getItem('token');
      const response = await fetch('http://aarya.live:8083/follow-status', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ other_user_id: userId }),
      });

      if (response.ok) {
        const data = await response.json();
        setFollowStatuses(prev => ({
          ...prev,
          [userId]: data,
        }));
      }
    } catch (error) {
      console.log('Error loading follow status:', error);
    }
  };

  const followUser = async (userId) => {
    try {
      const token = await AsyncStorage.getItem('token');
      const response = await fetch('http://aarya.live:8083/follow', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ other_user_id: userId }),
      });

      if (response.ok) {
        Alert.alert('Success', 'User followed successfully');
        loadFollowStatus(userId);
      } else {
        Alert.alert('Error', 'Failed to follow user');
      }
    } catch (error) {
      console.log('Error following user:', error);
      Alert.alert('Error', 'Failed to follow user');
    }
  };

  const unfollowUser = async (userId) => {
    try {
      const token = await AsyncStorage.getItem('token');
      const response = await fetch('http://aarya.live:8083/unfollow', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ other_user_id: userId }),
      });

      if (response.ok) {
        Alert.alert('Success', 'User unfollowed successfully');
        loadFollowStatus(userId);
      } else {
        Alert.alert('Error', 'Failed to unfollow user');
      }
    } catch (error) {
      console.log('Error unfollowing user:', error);
      Alert.alert('Error', 'Failed to unfollow user');
    }
  };

  const loadChatWithUser = (userId) => {
    if (!ws || !connected) {
      Alert.alert('Error', 'WebSocket not connected');
      return;
    }

    setLoading(true);
    const loadData = {
      type: 'load_chat',
      with: userId,
    };

    ws.send(JSON.stringify(loadData));
    setLoading(false);
  };

  const sendMessage = () => {
    if (!newMessage.trim() || !activeChat || !ws || !connected) {
      Alert.alert('Error', 'Please enter a message');
      return;
    }

    const messageData = {
      type: 'message',
      receiver_id: activeChat.id,
      content: newMessage.trim(),
    };

    ws.send(JSON.stringify(messageData));

    // Add message to local state optimistically
    const sentMsg = {
      sender_id: userInfo.id,
      receiver_id: activeChat.id,
      content: newMessage.trim(),
      created_at: new Date().toISOString(),
      is_read: false,
    };
    setMessages(prev => [...prev, sentMsg]);
    setNewMessage('');
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadChatList();
    setRefreshing(false);
  };

  const formatTime = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const formatDate = (timestamp) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
      return formatTime(timestamp);
    } else if (diffDays === 1) {
      return 'Yesterday';
    } else {
      return date.toLocaleDateString();
    }
  };

  const isMyMessage = (message) => {
    return message.sender_id === userInfo.id;
  };

  const renderHeader = () => (
    <View style={styles.header}>
      <TouchableOpacity 
        style={styles.backButton}
        onPress={() => {
          if (currentView === 'chat') {
            setCurrentView('chatList');
            setActiveChat(null);
            setMessages([]);
          } else {
            navigation.goBack();
          }
        }}
      >
        <Text style={styles.backButtonText}>← Back</Text>
      </TouchableOpacity>
      
      <Text style={styles.title}>
        {currentView === 'chatList' ? 'Chats' : 
         currentView === 'userSearch' ? 'Find Users' : 
         activeChat?.name || 'Chat'}
      </Text>
      
      <View style={styles.connectionStatus}>
        <View style={[styles.statusDot, { backgroundColor: connected ? '#4CAF50' : '#F44336' }]} />
        <Text style={styles.statusText}>{connected ? 'Connected' : 'Disconnected'}</Text>
      </View>
    </View>
  );

  const renderChatListItem = ({ item }) => (
    <TouchableOpacity
      style={styles.chatListItem}
      onPress={() => loadChatWithUser(item.user_id)}
    >
      <View style={styles.chatListItemContent}>
        <Text style={styles.chatListName}>{item.name}</Text>
        <Text style={styles.chatListEmail}>{item.email}</Text>
      </View>
      <View style={styles.chatListMeta}>
        <Text style={styles.chatListTime}>{formatDate(item.last_message_time)}</Text>
        {item.unread_count > 0 && (
          <View style={styles.unreadBadge}>
            <Text style={styles.unreadCount}>{item.unread_count}</Text>
          </View>
        )}
      </View>
    </TouchableOpacity>
  );

  const renderUserItem = ({ item }) => {
    const followStatus = followStatuses[item.id];
    
    return (
      <View style={styles.userItem}>
        <View style={styles.userItemContent}>
          <Text style={styles.userName}>{item.name}</Text>
          <Text style={styles.userEmail}>{item.email}</Text>
        </View>
        <View style={styles.userItemActions}>
          {followStatus?.mutual_follow && (
            <TouchableOpacity
              style={styles.chatButton}
              onPress={() => loadChatWithUser(item.id)}
            >
              <Text style={styles.chatButtonText}>Chat</Text>
            </TouchableOpacity>
          )}
          {followStatus?.i_follow_them ? (
            <TouchableOpacity
              style={styles.unfollowButton}
              onPress={() => unfollowUser(item.id)}
            >
              <Text style={styles.unfollowButtonText}>Unfollow</Text>
            </TouchableOpacity>
          ) : (
            <TouchableOpacity
              style={styles.followButton}
              onPress={() => followUser(item.id)}
            >
              <Text style={styles.followButtonText}>Follow</Text>
            </TouchableOpacity>
          )}
        </View>
      </View>
    );
  };

  const renderChatList = () => (
    <View style={styles.container}>
      {renderHeader()}
      
      <View style={styles.tabContainer}>
        <TouchableOpacity
          style={[styles.tab, currentView === 'chatList' && styles.activeTab]}
          onPress={() => setCurrentView('chatList')}
        >
          <Text style={[styles.tabText, currentView === 'chatList' && styles.activeTabText]}>
            Chats
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.tab, currentView === 'userSearch' && styles.activeTab]}
          onPress={() => {
            setCurrentView('userSearch');
            if (allUsers.length === 0) {
              loadAllUsers();
            }
          }}
        >
          <Text style={[styles.tabText, currentView === 'userSearch' && styles.activeTabText]}>
            Find Users
          </Text>
        </TouchableOpacity>
      </View>

      {currentView === 'chatList' ? (
        <FlatList
          data={chatList}
          renderItem={renderChatListItem}
          keyExtractor={(item) => item.user_id.toString()}
          style={styles.list}
          refreshControl={
            <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
          }
          ListEmptyComponent={
            <Text style={styles.emptyText}>No chats yet. Start a conversation!</Text>
          }
        />
      ) : (
        <View style={styles.userSearchContainer}>
          <TextInput
            style={styles.searchInput}
            placeholder="Search users by name or email"
            value={searchQuery}
            onChangeText={(text) => {
              setSearchQuery(text);
              searchUsers(text);
            }}
          />
          <FlatList
            data={searchQuery ? searchResults : allUsers}
            renderItem={renderUserItem}
            keyExtractor={(item) => item.id.toString()}
            style={styles.list}
            ListEmptyComponent={
              <Text style={styles.emptyText}>
                {searchQuery ? 'No users found' : 'No users available'}
              </Text>
            }
          />
        </View>
      )}
    </View>
  );

  const renderChat = () => (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView 
        style={styles.container} 
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      >
        {renderHeader()}

        {/* Messages */}
        <ScrollView
          ref={scrollViewRef}
          style={styles.messagesContainer}
          onContentSizeChange={() => scrollViewRef.current?.scrollToEnd()}
        >
          {loading && (
            <ActivityIndicator size="large" color="#007AFF" style={styles.loader} />
          )}
          
          {messages.map((message, index) => (
            <View
              key={index}
              style={[
                styles.messageBubble,
                isMyMessage(message) ? styles.myMessage : styles.theirMessage,
              ]}
            >
              <Text style={[
                styles.messageText,
                isMyMessage(message) ? styles.myMessageText : styles.theirMessageText
              ]}>
                {message.content}
              </Text>
              <Text style={[
                styles.messageTime,
                isMyMessage(message) ? styles.myMessageTime : styles.theirMessageTime
              ]}>
                {formatTime(message.created_at)}
                {!message.is_read && isMyMessage(message) && (
                  <Text style={styles.unreadIndicator}> • Unread</Text>
                )}
              </Text>
            </View>
          ))}
        </ScrollView>

        {/* Message Input */}
        <View style={styles.inputContainer}>
          <TextInput
            style={styles.messageInput}
            placeholder="Type a message..."
            value={newMessage}
            onChangeText={setNewMessage}
            multiline
            maxLength={500}
          />
          <TouchableOpacity
            style={[styles.sendButton, { opacity: newMessage.trim() ? 1 : 0.5 }]}
            onPress={sendMessage}
            disabled={!newMessage.trim()}
          >
            <Text style={styles.sendButtonText}>Send</Text>
          </TouchableOpacity>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );

  return currentView === 'chat' ? renderChat() : renderChatList();
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 15,
    backgroundColor: 'white',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  backButton: {
    padding: 5,
  },
  backButtonText: {
    fontSize: 16,
    color: '#007AFF',
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#333',
  },
  connectionStatus: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 5,
  },
  statusText: {
    fontSize: 12,
    color: '#666',
  },
  tabContainer: {
    flexDirection: 'row',
    backgroundColor: 'white',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  tab: {
    flex: 1,
    paddingVertical: 15,
    alignItems: 'center',
    borderBottomWidth: 2,
    borderBottomColor: 'transparent',
  },
  activeTab: {
    borderBottomColor: '#007AFF',
  },
  tabText: {
    fontSize: 16,
    color: '#666',
  },
  activeTabText: {
    color: '#007AFF',
    fontWeight: 'bold',
  },
  list: {
    flex: 1,
  },
  chatListItem: {
    flexDirection: 'row',
    padding: 15,
    backgroundColor: 'white',
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
    alignItems: 'center',
  },
  chatListItemContent: {
    flex: 1,
  },
  chatListName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#333',
  },
  chatListEmail: {
    fontSize: 14,
    color: '#666',
    marginTop: 2,
  },
  chatListMeta: {
    alignItems: 'flex-end',
  },
  chatListTime: {
    fontSize: 12,
    color: '#666',
  },
  unreadBadge: {
    backgroundColor: '#007AFF',
    borderRadius: 10,
    paddingHorizontal: 6,
    paddingVertical: 2,
    marginTop: 4,
    minWidth: 20,
    alignItems: 'center',
  },
  unreadCount: {
    color: 'white',
    fontSize: 12,
    fontWeight: 'bold',
  },
  userSearchContainer: {
    flex: 1,
  },
  searchInput: {
    margin: 15,
    padding: 12,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    fontSize: 16,
    backgroundColor: 'white',
  },
  userItem: {
    flexDirection: 'row',
    padding: 15,
    backgroundColor: 'white',
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
    alignItems: 'center',
  },
  userItemContent: {
    flex: 1,
  },
  userName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#333',
  },
  userEmail: {
    fontSize: 14,
    color: '#666',
    marginTop: 2,
  },
  userItemActions: {
    flexDirection: 'row',
    gap: 10,
  },
  followButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  followButtonText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 14,
  },
  unfollowButton: {
    backgroundColor: '#FF3B30',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  unfollowButtonText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 14,
  },
  chatButton: {
    backgroundColor: '#34C759',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  chatButtonText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 14,
  },
  emptyText: {
    textAlign: 'center',
    color: '#666',
    fontSize: 16,
    marginTop: 50,
  },
  messagesContainer: {
    flex: 1,
    padding: 15,
  },
  loader: {
    marginVertical: 20,
  },
  messageBubble: {
    maxWidth: '80%',
    padding: 12,
    borderRadius: 18,
    marginBottom: 10,
  },
  myMessage: {
    alignSelf: 'flex-end',
    backgroundColor: '#007AFF',
  },
  theirMessage: {
    alignSelf: 'flex-start',
    backgroundColor: '#E5E5EA',
  },
  messageText: {
    fontSize: 16,
  },
  myMessageText: {
    color: 'white',
  },
  theirMessageText: {
    color: '#333',
  },
  messageTime: {
    fontSize: 12,
    marginTop: 4,
  },
  myMessageTime: {
    color: 'rgba(255, 255, 255, 0.7)',
  },
  theirMessageTime: {
    color: '#666',
  },
  unreadIndicator: {
    color: '#FFD60A',
  },
  inputContainer: {
    flexDirection: 'row',
    padding: 15,
    backgroundColor: 'white',
    alignItems: 'flex-end',
  },
  messageInput: {
    flex: 1,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 20,
    paddingHorizontal: 15,
    paddingVertical: 10,
    marginRight: 10,
    fontSize: 16,
    maxHeight: 100,
  },
  sendButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 20,
  },
  sendButtonText: {
    color: 'white',
    fontWeight: 'bold',
  },
});