// // screens/VideoScreen.js
// import React, { useState, useEffect } from 'react';
// import {
//   View,
//   Text,
//   StyleSheet,
//   ScrollView,
//   TouchableOpacity,
//   Alert,
//   ActivityIndicator,
//   RefreshControl,
//   Dimensions,
//   Platform,
// } from 'react-native';
// import * as ImagePicker from 'expo-image-picker';
// import { Video } from 'expo-av';
// import AsyncStorage from '@react-native-async-storage/async-storage';
// import { useFocusEffect } from '@react-navigation/native';
// import { MaterialIcons } from '@expo/vector-icons';

// const { width } = Dimensions.get('window');

// const VideoScreen = ({ navigation, userInfo }) => {
//   const [videos, setVideos] = useState([]);
//   const [loading, setLoading] = useState(false);
//   const [refreshing, setRefreshing] = useState(false);
//   const [uploading, setUploading] = useState(false);
//   const [playingVideo, setPlayingVideo] = useState(null);

//   // Your backend URL - replace with your actual URL
//   const API_BASE_URL = 'http://10.0.2.2:8084';

//   // Request permissions
//   const requestPermissions = async () => {
//     const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
//     if (status !== 'granted') {
//       Alert.alert('Permission needed', 'Please grant media library permissions to upload videos');
//       return false;
//     }
//     return true;
//   };

//   // Fetch user's videos
//   const fetchVideos = async () => {
//     try {
//       setLoading(true);
//       const token = await AsyncStorage.getItem('token');
      
//       if (!token) {
//         Alert.alert('Error', 'Please login first');
//         navigation.navigate('Login');
//         return;
//       }

//       console.log('Fetching videos with token:', token.substring(0, 20) + '...');
      
//       const response = await fetch(`${API_BASE_URL}/videos`, {
//         method: 'GET',
//         headers: {
//           'Authorization': `Bearer ${token}`,
//           'Content-Type': 'application/json',
//         },
//       });

//       console.log('Response status:', response.status);
//       console.log('Response headers:', response.headers);

//       if (response.ok) {
//         const data = await response.json();
//         console.log('Videos fetched:', data);
//         setVideos(data || []);
//       } else if (response.status === 401) {
//         console.log('Unauthorized - token expired');
//         Alert.alert('Session Expired', 'Please login again');
//         navigation.navigate('Home');
//       } else {
//         const errorText = await response.text();
//         console.log('Error response:', errorText);
//         throw new Error(`HTTP ${response.status}: ${errorText}`);
//       }
//     } catch (error) {
//       console.error('Error fetching videos:', error);
//       Alert.alert('Error', `Failed to load videos: ${error.message}`);
//     } finally {
//       setLoading(false);
//     }
//   };

//   // Upload video to S3 using presigned URL
//   const uploadVideoToS3 = async (uploadUrl, videoUri) => {
//     try {
//       console.log('Uploading video to S3:', uploadUrl);
      
//       // For React Native, we need to use FormData with the file
//       const formData = new FormData();
//       formData.append('file', {
//         uri: videoUri,
//         type: 'video/mp4',
//         name: 'video.mp4',
//       });

//       const response = await fetch(uploadUrl, {
//         method: 'PUT',
//         body: formData,
//         headers: {
//           'Content-Type': 'multipart/form-data',
//         },
//       });

//       console.log('S3 upload response status:', response.status);
//       return response.ok;
//     } catch (error) {
//       console.error('S3 upload error:', error);
//       return false;
//     }
//   };

//   // Alternative S3 upload method using raw file data
// //   const uploadVideoToS3Alternative = async (uploadUrl, videoUri) => {
// //     try {
// //       console.log('Trying alternative S3 upload method');
      
// //       // Read file as blob/binary data
// //       const fileResponse = await fetch(videoUri);
// //       const blob = await fileResponse.blob();

// //       const response = await fetch(uploadUrl, {
// //         method: 'PUT',
// //         body: blob,
// //         headers: {
// //           'Content-Type': 'video/mp4',
// //         },
// //       });

// //       console.log('Alternative S3 upload response status:', response.status);
// //       return response.ok;
// //     } catch (error) {
// //       console.error('Alternative S3 upload error:', error);
// //       return false;
// //     }
// //   };

// //     const uploadVideoToS3Alternative = async (uploadUrl, videoUri) => {
// //         try {
// //             console.log('Trying alternative S3 upload method');

// //             const fileResponse = await fetch(videoUri);
// //             const blob = await fileResponse.blob();

// //             const response = await fetch(uploadUrl, {
// //             method: 'PUT',
// //             body: blob,
// //             headers: {
// //                 'Content-Type': 'video/mp4', // must match exactly
// //                 'x-amz-acl': 'public-read',  // include this if used during presign
// //             },
// //             });

// //             console.log('Alternative S3 upload response status:', response.status);
// //             return response.ok;
// //         } catch (error) {
// //             console.error('Alternative S3 upload error:', error);
// //             return false;
// //         }
// //         };


// //   // Handle video selection and upload
// //   const handleVideoUpload = async () => {
// //     console.log('🎬 Video upload button pressed!');
    
// //     try {
// //       // Request permissions first
// //       const hasPermission = await requestPermissions();
// //       if (!hasPermission) {
// //         console.log('❌ Permission denied');
// //         return;
// //       }

// //       console.log('✅ Permission granted, launching image picker...');

// //       // Launch image picker for video
// //       const result = await ImagePicker.launchImageLibraryAsync({
// //         mediaTypes: ImagePicker.MediaTypeOptions.Videos,
// //         allowsEditing: true,
// //         quality: 0.8,
// //         duration: 300000, // 5 minutes max (in milliseconds)
// //       });

// //       console.log('📱 Image picker result:', result);

// //       if (result.canceled) {
// //         console.log('User cancelled video selection');
// //         return;
// //       }

// //       if (!result.assets || result.assets.length === 0) {
// //         Alert.alert('Error', 'No video selected');
// //         return;
// //       }

// //       const video = result.assets[0];
// //       console.log('Selected video:', video);

// //       try {
// //         setUploading(true);
// //         const token = await AsyncStorage.getItem('token');

// //         if (!token) {
// //           Alert.alert('Error', 'Please login first');
// //           navigation.navigate('Home');
// //           return;
// //         }

// //         // Step 1: Get presigned upload URL
// //         console.log('Getting presigned upload URL...');
// //         const uploadResponse = await fetch(`${API_BASE_URL}/generate-upload-url`, {
// //           method: 'POST',
// //           headers: {
// //             'Authorization': `Bearer ${token}`,
// //             'Content-Type': 'application/json',
// //           },
// //           body: JSON.stringify({
// //             filename: video.fileName || `video_${Date.now()}.mp4`,
// //             contentType: 'video/mp4',
// //           }),
// //         });

// //         console.log('Upload URL response status:', uploadResponse.status);

// //         if (!uploadResponse.ok) {
// //           const errorText = await uploadResponse.text();
// //           console.log('Upload URL error:', errorText);
// //           throw new Error(`Failed to get upload URL: ${errorText}`);
// //         }

// //         const { upload_url, video_id } = await uploadResponse.json();
// //         console.log('Got upload URL and video ID:', { upload_url: upload_url.substring(0, 50) + '...', video_id });

// //         // Step 2: Upload video to S3
// //         console.log('Uploading video to S3...');
// //         let uploadSuccess = await uploadVideoToS3Alternative(upload_url, video.uri);
        
// //         if (!uploadSuccess) {
// //           console.log('Alternative upload failed, trying FormData method...');
// //           uploadSuccess = await uploadVideoToS3(upload_url, video.uri);
// //         }

// //         if (uploadSuccess) {
// //           console.log('Video uploaded successfully!');
// //           Alert.alert('Success', 'Video uploaded successfully!');
// //           fetchVideos(); // Refresh the video list
// //         } else {
// //           throw new Error('Failed to upload video to S3');
// //         }
// //       } catch (error) {
// //         console.error('Upload error:', error);
// //         Alert.alert('Error', error.message || 'Failed to upload video');
// //       } finally {
// //         setUploading(false);
// //       }
// //     } catch (error) {
// //       console.error('Video selection error:', error);
// //       Alert.alert('Error', 'Failed to select video');
// //     }
// //   };



// // Fixed uploadVideoToS3Alternative function
// const uploadVideoToS3Alternative = async (uploadUrl, videoUri) => {
//   try {
//     console.log('Uploading video to S3 with binary data');
    
//     // Read the file as binary data
//     const response = await fetch(videoUri);
//     const blob = await response.blob();
    
//     console.log('File blob created, size:', blob.size);
    
//     // Upload to S3 with exact same headers as your working Postman request
//     const uploadResponse = await fetch(uploadUrl, {
//       method: 'PUT',
//       body: blob,
//       headers: {
//         'Content-Type': 'video/mp4', // Must match exactly what was signed
//       },
//     });
    
//     console.log('S3 upload response status:', uploadResponse.status);
    
//     if (!uploadResponse.ok) {
//       const errorText = await uploadResponse.text();
//       console.log('S3 upload error response:', errorText);
//       throw new Error(`S3 upload failed: ${uploadResponse.status} ${errorText}`);
//     }
    
//     return true;
//   } catch (error) {
//     console.error('S3 upload error:', error);
//     throw error;
//   }
// };

// // Remove the FormData upload method entirely and update the main upload function
// const handleVideoUpload = async () => {
//   console.log('🎬 Video upload button pressed!');
  
//   try {
//     // Request permissions first
//     const hasPermission = await requestPermissions();
//     if (!hasPermission) {
//       console.log('❌ Permission denied');
//       return;
//     }

//     console.log('✅ Permission granted, launching image picker...');

//     // Launch image picker for video
//     const result = await ImagePicker.launchImageLibraryAsync({
//       mediaTypes: ImagePicker.MediaType.Videos, // Fixed deprecated MediaTypeOptions
//       allowsEditing: true,
//       quality: 0.8,
//       duration: 300000, // 5 minutes max
//     });

//     console.log('📱 Image picker result:', result);

//     if (result.canceled) {
//       console.log('User cancelled video selection');
//       return;
//     }

//     if (!result.assets || result.assets.length === 0) {
//       Alert.alert('Error', 'No video selected');
//       return;
//     }

//     const video = result.assets[0];
//     console.log('Selected video:', video);

//     try {
//       setUploading(true);
//       const token = await AsyncStorage.getItem('token');

//       if (!token) {
//         Alert.alert('Error', 'Please login first');
//         navigation.navigate('Home');
//         return;
//       }

//       // Step 1: Get presigned upload URL
//       console.log('Getting presigned upload URL...');
//       const uploadResponse = await fetch(`${API_BASE_URL}/generate-upload-url`, {
//         method: 'POST',
//         headers: {
//           'Authorization': `Bearer ${token}`,
//           'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({
//           filename: video.fileName || `video_${Date.now()}.mp4`,
//         }),
//       });

//       console.log('Upload URL response status:', uploadResponse.status);

//       if (!uploadResponse.ok) {
//         const errorText = await uploadResponse.text();
//         console.log('Upload URL error:', errorText);
//         throw new Error(`Failed to get upload URL: ${errorText}`);
//       }

//       const { upload_url, video_id } = await uploadResponse.json();
//       console.log('Got upload URL and video ID:', { 
//         upload_url: upload_url.substring(0, 50) + '...', 
//         video_id 
//       });

//       // Step 2: Upload video to S3 using the fixed method
//       console.log('Uploading video to S3...');
//       await uploadVideoToS3Alternative(upload_url, video.uri);
      
//       console.log('Video uploaded successfully!');
//       Alert.alert('Success', 'Video uploaded successfully!');
//       fetchVideos(); // Refresh the video list

//     } catch (error) {
//       console.error('Upload error:', error);
//       Alert.alert('Error', error.message || 'Failed to upload video');
//     } finally {
//       setUploading(false);
//     }
//   } catch (error) {
//     console.error('Video selection error:', error);
//     Alert.alert('Error', 'Failed to select video');
//   }
// };

//   // Handle video play/pause
//   const toggleVideoPlay = (videoId) => {
//     setPlayingVideo(playingVideo === videoId ? null : videoId);
//   };

//   // Handle delete video
//   const handleDeleteVideo = async (videoId) => {
//     Alert.alert(
//       'Delete Video',
//       'Are you sure you want to delete this video?',
//       [
//         { text: 'Cancel', style: 'cancel' },
//         { 
//           text: 'Delete', 
//           style: 'destructive',
//           onPress: async () => {
//             try {
//               const token = await AsyncStorage.getItem('token');
//               if (!token) {
//                 Alert.alert('Error', 'Please login first');
//                 return;
//               }

//               console.log('Deleting video:', videoId);
//               const response = await fetch(`${API_BASE_URL}/videos/${videoId}`, {
//                 method: 'DELETE',
//                 headers: {
//                   'Authorization': `Bearer ${token}`,
//                 },
//               });

//               console.log('Delete response status:', response.status);

//               if (response.ok) {
//                 Alert.alert('Success', 'Video deleted successfully');
//                 fetchVideos(); // Refresh the list
//               } else {
//                 const errorText = await response.text();
//                 console.log('Delete error:', errorText);
//                 throw new Error(`Failed to delete video: ${errorText}`);
//               }
//             } catch (error) {
//               console.error('Delete error:', error);
//               Alert.alert('Error', error.message || 'Failed to delete video');
//             }
//           }
//         }
//       ]
//     );
//   };

//   // Refresh videos
//   const onRefresh = async () => {
//     setRefreshing(true);
//     await fetchVideos();
//     setRefreshing(false);
//   };

//   // Load videos when screen is focused
//   useFocusEffect(
//     React.useCallback(() => {
//       fetchVideos();
//     }, [])
//   );

//   // Render video item
//   const renderVideoItem = (video) => (
//     <View key={video.id} style={styles.videoCard}>
//       <View style={styles.videoHeader}>
//         <Text style={styles.videoTitle} numberOfLines={1}>
//           {video.video_name?.replace(/^\d+_/, '') || 'Untitled Video'}
//         </Text>
//         <TouchableOpacity
//           style={styles.deleteButton}
//           onPress={() => handleDeleteVideo(video.id)}
//         >
//           <MaterialIcons name="delete" size={24} color="#ff4444" />
//         </TouchableOpacity>
//       </View>
      
//       <View style={styles.videoContainer}>
//         {playingVideo === video.id ? (
//           <Video
//             source={{ uri: video.s3_url }}
//             style={styles.videoPlayer}
//             useNativeControls
//             resizeMode="contain"
//             onError={(error) => {
//               console.error('Video playback error:', error);
//               Alert.alert('Error', 'Failed to play video');
//               setPlayingVideo(null);
//             }}
//             onPlaybackStatusUpdate={(status) => {
//               if (status.didJustFinish) {
//                 setPlayingVideo(null);
//               }
//             }}
//           />
//         ) : (
//           <TouchableOpacity
//             style={styles.videoThumbnail}
//             onPress={() => toggleVideoPlay(video.id)}
//           >
//             <MaterialIcons name="play-circle-filled" size={64} color="#007bff" />
//             <Text style={styles.playText}>Tap to play</Text>
//           </TouchableOpacity>
//         )}
//       </View>
      
//       <View style={styles.videoInfo}>
//         <Text style={styles.videoDate}>
//           {new Date(video.timestamp).toLocaleDateString()}
//         </Text>
//         {playingVideo === video.id && (
//           <TouchableOpacity
//             style={styles.stopButton}
//             onPress={() => setPlayingVideo(null)}
//           >
//             <MaterialIcons name="stop" size={20} color="#666" />
//             <Text style={styles.stopText}>Stop</Text>
//           </TouchableOpacity>
//         )}
//       </View>
//     </View>
//   );

//   return (
//     <View style={styles.container}>
//       {/* Header */}
//       <View style={styles.header}>
//         <TouchableOpacity
//           style={styles.backButton}
//           onPress={() => navigation.goBack()}
//         >
//           <MaterialIcons name="arrow-back" size={24} color="#333" />
//         </TouchableOpacity>
//         <Text style={styles.headerTitle}>My Videos</Text>
//         <TouchableOpacity
//           style={styles.uploadButton}
//           onPress={handleVideoUpload}
//           disabled={uploading}
//         >
//           {uploading ? (
//             <ActivityIndicator color="#007bff" />
//           ) : (
//             <MaterialIcons name="add" size={24} color="#007bff" />
//           )}
//         </TouchableOpacity>
//       </View>

//       {/* User Info */}
//       {userInfo && (
//         <View style={styles.userInfo}>
//           <Text style={styles.userInfoText}>
//             Welcome, {userInfo.name}! 👋
//           </Text>
//         </View>
//       )}

//       {/* Content */}
//       <ScrollView
//         style={styles.content}
//         refreshControl={
//           <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
//         }
//         showsVerticalScrollIndicator={false}
//       >
//         {loading ? (
//           <View style={styles.loadingContainer}>
//             <ActivityIndicator size="large" color="#007bff" />
//             <Text style={styles.loadingText}>Loading videos...</Text>
//           </View>
//         ) : videos.length === 0 ? (
//           <View style={styles.emptyContainer}>
//             <MaterialIcons name="video-library" size={64} color="#ccc" />
//             <Text style={styles.emptyText}>No videos yet</Text>
//             <Text style={styles.emptySubtext}>
//               Tap the + button to upload your first video
//             </Text>
//           </View>
//         ) : (
//           videos.map(renderVideoItem)
//         )}
//       </ScrollView>

//       {/* Upload Progress */}
//       {uploading && (
//         <View style={styles.uploadOverlay}>
//           <View style={styles.uploadProgress}>
//             <ActivityIndicator size="large" color="#007bff" />
//             <Text style={styles.uploadText}>Uploading video...</Text>
//             <Text style={styles.uploadSubtext}>Please wait...</Text>
//           </View>
//         </View>
//       )}
//     </View>
//   );
// };

// const styles = StyleSheet.create({
//   container: {
//     flex: 1,
//     backgroundColor: '#f5f5f5',
//   },
//   header: {
//     flexDirection: 'row',
//     alignItems: 'center',
//     justifyContent: 'space-between',
//     padding: 16,
//     backgroundColor: '#fff',
//     borderBottomWidth: 1,
//     borderBottomColor: '#e0e0e0',
//     elevation: 2,
//     shadowColor: '#000',
//     shadowOffset: { width: 0, height: 1 },
//     shadowOpacity: 0.2,
//     shadowRadius: 2,
//   },
//   backButton: {
//     padding: 8,
//   },
//   headerTitle: {
//     fontSize: 20,
//     fontWeight: 'bold',
//     color: '#333',
//   },
//   uploadButton: {
//     padding: 8,
//   },
//   userInfo: {
//     backgroundColor: '#fff',
//     padding: 16,
//     borderBottomWidth: 1,
//     borderBottomColor: '#e0e0e0',
//   },
//   userInfoText: {
//     fontSize: 16,
//     color: '#333',
//     textAlign: 'center',
//   },
//   content: {
//     flex: 1,
//     padding: 16,
//   },
//   loadingContainer: {
//     flex: 1,
//     justifyContent: 'center',
//     alignItems: 'center',
//     paddingVertical: 50,
//   },
//   loadingText: {
//     marginTop: 10,
//     fontSize: 16,
//     color: '#666',
//   },
//   emptyContainer: {
//     flex: 1,
//     justifyContent: 'center',
//     alignItems: 'center',
//     paddingVertical: 50,
//   },
//   emptyText: {
//     fontSize: 18,
//     color: '#666',
//     marginTop: 16,
//   },
//   emptySubtext: {
//     fontSize: 14,
//     color: '#999',
//     textAlign: 'center',
//     marginTop: 8,
//   },
//   videoCard: {
//     backgroundColor: '#fff',
//     borderRadius: 12,
//     padding: 16,
//     marginBottom: 16,
//     elevation: 3,
//     shadowColor: '#000',
//     shadowOffset: { width: 0, height: 2 },
//     shadowOpacity: 0.1,
//     shadowRadius: 4,
//   },
//   videoHeader: {
//     flexDirection: 'row',
//     justifyContent: 'space-between',
//     alignItems: 'center',
//     marginBottom: 12,
//   },
//   videoTitle: {
//     fontSize: 16,
//     fontWeight: 'bold',
//     color: '#333',
//     flex: 1,
//   },
//   deleteButton: {
//     padding: 8,
//   },
//   videoContainer: {
//     width: '100%',
//     height: 200,
//     borderRadius: 8,
//     overflow: 'hidden',
//     backgroundColor: '#000',
//   },
//   videoPlayer: {
//     width: '100%',
//     height: '100%',
//   },
//   videoThumbnail: {
//     width: '100%',
//     height: '100%',
//     justifyContent: 'center',
//     alignItems: 'center',
//     backgroundColor: '#f0f0f0',
//   },
//   playText: {
//     marginTop: 8,
//     fontSize: 14,
//     color: '#007bff',
//   },
//   videoInfo: {
//     flexDirection: 'row',
//     justifyContent: 'space-between',
//     alignItems: 'center',
//     marginTop: 12,
//   },
//   videoDate: {
//     fontSize: 12,
//     color: '#666',
//   },
//   stopButton: {
//     flexDirection: 'row',
//     alignItems: 'center',
//     padding: 4,
//   },
//   stopText: {
//     marginLeft: 4,
//     fontSize: 12,
//     color: '#666',
//   },
//   uploadOverlay: {
//     position: 'absolute',
//     top: 0,
//     left: 0,
//     right: 0,
//     bottom: 0,
//     backgroundColor: 'rgba(0, 0, 0, 0.5)',
//     justifyContent: 'center',
//     alignItems: 'center',
//   },
//   uploadProgress: {
//     backgroundColor: '#fff',
//     padding: 20,
//     borderRadius: 10,
//     alignItems: 'center',
//     minWidth: 200,
//   },
//   uploadText: {
//     marginTop: 10,
//     fontSize: 16,
//     color: '#333',
//   },
//   uploadSubtext: {
//     marginTop: 5,
//     fontSize: 14,
//     color: '#666',
//   },
// });

// export default VideoScreen;

// screens/VideoScreen.js
import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
  ActivityIndicator,
  RefreshControl,
  Dimensions,
  Platform,
} from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import * as FileSystem from 'expo-file-system';
import { Video } from 'expo-av';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { useFocusEffect } from '@react-navigation/native';
import { MaterialIcons } from '@expo/vector-icons';

const { width } = Dimensions.get('window');

const VideoScreen = ({ navigation, userInfo }) => {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [playingVideo, setPlayingVideo] = useState(null);

  // Your backend URL - replace with your actual URL
  const API_BASE_URL = 'http://aarya.live:8084';

  // Request permissions
  const requestPermissions = async () => {
    const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (status !== 'granted') {
      Alert.alert('Permission needed', 'Please grant media library permissions to upload videos');
      return false;
    }
    return true;
  };

  // Fetch user's videos
  const fetchVideos = async () => {
    try {
      setLoading(true);
      const token = await AsyncStorage.getItem('token');
      
      if (!token) {
        Alert.alert('Error', 'Please login first');
        navigation.navigate('Login');
        return;
      }

      console.log('Fetching videos with token:', token.substring(0, 20) + '...');
      
      const response = await fetch(`${API_BASE_URL}/videos`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      console.log('Response status:', response.status);
      console.log('Response headers:', response.headers);

      if (response.ok) {
        const data = await response.json();
        console.log('Videos fetched:', data);
        setVideos(data || []);
      } else if (response.status === 401) {
        console.log('Unauthorized - token expired');
        Alert.alert('Session Expired', 'Please login again');
        navigation.navigate('Home');
      } else {
        const errorText = await response.text();
        console.log('Error response:', errorText);
        throw new Error(`HTTP ${response.status}: ${errorText}`);
      }
    } catch (error) {
      console.error('Error fetching videos:', error);
      Alert.alert('Error', `Failed to load videos: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const uploadVideoToS3 = async (uploadUrl, videoUri) => {
  try {
    console.log('Uploading video to S3 using FileSystem');
    console.log('Video URI:', videoUri);
    console.log('Upload URL:', uploadUrl.substring(0, 50) + '...');
    
    // Get file info first
    const fileInfo = await FileSystem.getInfoAsync(videoUri);
    console.log('File info:', fileInfo);
    
    if (!fileInfo.exists) {
      throw new Error('File does not exist at path: ' + videoUri);
    }
    
    // Use FileSystem.uploadAsync for reliable binary uploads
    const uploadResponse = await FileSystem.uploadAsync(uploadUrl, fileInfo.uri, {
      httpMethod: 'PUT',
      uploadType: FileSystem.FileSystemUploadType.BINARY_CONTENT,
      headers: {
        'Content-Type': 'video/mp4',
      },
    });
    
    console.log('FileSystem upload response:', uploadResponse);
    
    if (uploadResponse.status !== 200) {
      throw new Error(`S3 upload failed: ${uploadResponse.status}`);
    }
    
    return true;
  } catch (error) {
    console.error('S3 upload error:', error);
    throw error;
  }
};

  // Handle video selection and upload
  const handleVideoUpload = async () => {
    console.log('🎬 Video upload button pressed!');
    
    try {
      // Request permissions first
      const hasPermission = await requestPermissions();
      if (!hasPermission) {
        console.log('❌ Permission denied');
        return;
      }

      console.log('✅ Permission granted, launching image picker...');

      // Launch image picker for video - Fixed deprecated API
      const result = await ImagePicker.launchImageLibraryAsync({
        mediaTypes: ['videos'], // Use array instead of deprecated MediaTypeOptions
        allowsEditing: true,
        quality: 0.8,
        duration: 300000, // 5 minutes max (in milliseconds)
      });

      console.log('📱 Image picker result:', result);

      if (result.canceled) {
        console.log('User cancelled video selection');
        return;
      }

      if (!result.assets || result.assets.length === 0) {
        Alert.alert('Error', 'No video selected');
        return;
      }

      const video = result.assets[0];
      console.log('Selected video:', video);

      try {
        setUploading(true);
        const token = await AsyncStorage.getItem('token');

        if (!token) {
          Alert.alert('Error', 'Please login first');
          navigation.navigate('Home');
          return;
        }

        // Step 1: Get presigned upload URL
        console.log('Getting presigned upload URL...');
        const uploadResponse = await fetch(`${API_BASE_URL}/generate-upload-url`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            filename: video.fileName || `video_${Date.now()}.mp4`,
          }),
        });

        console.log('Upload URL response status:', uploadResponse.status);

        if (!uploadResponse.ok) {
          const errorText = await uploadResponse.text();
          console.log('Upload URL error:', errorText);
          throw new Error(`Failed to get upload URL: ${errorText}`);
        }

        const { upload_url, video_id } = await uploadResponse.json();
        console.log('Got upload URL and video ID:');
            // , { 
        //   upload_url: upload_url.substring(0, 50) + '...', 
        //   video_id 
        // });

        // Step 2: Upload video to S3 using the fixed method
        console.log('Uploading video to S3...');
        await uploadVideoToS3(upload_url, video.uri);
        
        console.log('Video uploaded successfully!');
        Alert.alert('Success', 'Video uploaded successfully!');
        fetchVideos(); // Refresh the video list

      } catch (error) {
        console.error('Upload error:', error);
        Alert.alert('Error', error.message || 'Failed to upload video');
      } finally {
        setUploading(false);
      }
    } catch (error) {
      console.error('Video selection error:', error);
      Alert.alert('Error', 'Failed to select video');
    }
  };

  // Handle video play/pause
  const toggleVideoPlay = (videoId) => {
    setPlayingVideo(playingVideo === videoId ? null : videoId);
  };

  // Handle delete video
  const handleDeleteVideo = async (videoId) => {
    Alert.alert(
      'Delete Video',
      'Are you sure you want to delete this video?',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Delete', 
          style: 'destructive',
          onPress: async () => {
            try {
              const token = await AsyncStorage.getItem('token');
              if (!token) {
                Alert.alert('Error', 'Please login first');
                return;
              }

              console.log('Deleting video:', videoId);
              const response = await fetch(`${API_BASE_URL}/videos/${videoId}`, {
                method: 'DELETE',
                headers: {
                  'Authorization': `Bearer ${token}`,
                },
              });

              console.log('Delete response status:', response.status);

              if (response.ok) {
                Alert.alert('Success', 'Video deleted successfully');
                fetchVideos(); // Refresh the list
              } else {
                const errorText = await response.text();
                console.log('Delete error:', errorText);
                throw new Error(`Failed to delete video: ${errorText}`);
              }
            } catch (error) {
              console.error('Delete error:', error);
              Alert.alert('Error', error.message || 'Failed to delete video');
            }
          }
        }
      ]
    );
  };

  // Refresh videos
  const onRefresh = async () => {
    setRefreshing(true);
    await fetchVideos();
    setRefreshing(false);
  };

  // Load videos when screen is focused
  useFocusEffect(
    React.useCallback(() => {
      fetchVideos();
    }, [])
  );

  // Render video item
  const renderVideoItem = (video) => (
    <View key={video.id} style={styles.videoCard}>
      <View style={styles.videoHeader}>
        <Text style={styles.videoTitle} numberOfLines={1}>
          {video.video_name?.replace(/^\d+_/, '') || 'Untitled Video'}
        </Text>
        <TouchableOpacity
          style={styles.deleteButton}
          onPress={() => handleDeleteVideo(video.id)}
        >
          <MaterialIcons name="delete" size={24} color="#ff4444" />
        </TouchableOpacity>
      </View>
      
      <View style={styles.videoContainer}>
        {playingVideo === video.id ? (
          <Video
            source={{ uri: video.s3_url }}
            style={styles.videoPlayer}
            useNativeControls
            resizeMode="contain"
            onError={(error) => {
              console.error('Video playback error:', error);
              Alert.alert('Error', 'Failed to play video');
              setPlayingVideo(null);
            }}
            onPlaybackStatusUpdate={(status) => {
              if (status.didJustFinish) {
                setPlayingVideo(null);
              }
            }}
          />
        ) : (
          <TouchableOpacity
            style={styles.videoThumbnail}
            onPress={() => toggleVideoPlay(video.id)}
          >
            <MaterialIcons name="play-circle-filled" size={64} color="#007bff" />
            <Text style={styles.playText}>Tap to play</Text>
          </TouchableOpacity>
        )}
      </View>
      
      <View style={styles.videoInfo}>
        <Text style={styles.videoDate}>
          {new Date(video.timestamp).toLocaleDateString()}
        </Text>
        {playingVideo === video.id && (
          <TouchableOpacity
            style={styles.stopButton}
            onPress={() => setPlayingVideo(null)}
          >
            <MaterialIcons name="stop" size={20} color="#666" />
            <Text style={styles.stopText}>Stop</Text>
          </TouchableOpacity>
        )}
      </View>
    </View>
  );

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity
          style={styles.backButton}
          onPress={() => navigation.goBack()}
        >
          <MaterialIcons name="arrow-back" size={24} color="#333" />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>My Videos</Text>
        <TouchableOpacity
          style={styles.uploadButton}
          onPress={handleVideoUpload}
          disabled={uploading}
        >
          {uploading ? (
            <ActivityIndicator color="#007bff" />
          ) : (
            <MaterialIcons name="add" size={24} color="#007bff" />
          )}
        </TouchableOpacity>
      </View>

      {/* User Info */}
      {userInfo && (
        <View style={styles.userInfo}>
          <Text style={styles.userInfoText}>
            Welcome, {userInfo.name}! 👋
          </Text>
        </View>
      )}

      {/* Content */}
      <ScrollView
        style={styles.content}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
        showsVerticalScrollIndicator={false}
      >
        {loading ? (
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" color="#007bff" />
            <Text style={styles.loadingText}>Loading videos...</Text>
          </View>
        ) : videos.length === 0 ? (
          <View style={styles.emptyContainer}>
            <MaterialIcons name="video-library" size={64} color="#ccc" />
            <Text style={styles.emptyText}>No videos yet</Text>
            <Text style={styles.emptySubtext}>
              Tap the + button to upload your first video
            </Text>
          </View>
        ) : (
          videos.map(renderVideoItem)
        )}
      </ScrollView>

      {/* Upload Progress */}
      {uploading && (
        <View style={styles.uploadOverlay}>
          <View style={styles.uploadProgress}>
            <ActivityIndicator size="large" color="#007bff" />
            <Text style={styles.uploadText}>Uploading video...</Text>
            <Text style={styles.uploadSubtext}>Please wait...</Text>
          </View>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingTop: Platform.OS === 'ios' ? 50 : 30,
    paddingBottom: 16,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  backButton: {
    padding: 8,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  uploadButton: {
    padding: 8,
  },
  userInfo: {
    backgroundColor: '#fff',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  userInfoText: {
    fontSize: 16,
    color: '#333',
  },
  content: {
    flex: 1,
    padding: 16,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingTop: 50,
  },
  loadingText: {
    marginTop: 16,
    fontSize: 16,
    color: '#666',
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingTop: 50,
  },
  emptyText: {
    fontSize: 18,
    color: '#666',
    marginTop: 16,
  },
  emptySubtext: {
    fontSize: 14,
    color: '#999',
    textAlign: 'center',
    marginTop: 8,
  },
  videoCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  videoHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  videoTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    flex: 1,
    marginRight: 8,
  },
  deleteButton: {
    padding: 4,
  },
  videoContainer: {
    borderRadius: 8,
    overflow: 'hidden',
    backgroundColor: '#f0f0f0',
  },
  videoPlayer: {
    width: '100%',
    height: 200,
  },
  videoThumbnail: {
    width: '100%',
    height: 200,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f0f0f0',
  },
  playText: {
    marginTop: 8,
    fontSize: 14,
    color: '#007bff',
  },
  videoInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 12,
  },
  videoDate: {
    fontSize: 12,
    color: '#666',
  },
  stopButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f0f0f0',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  stopText: {
    marginLeft: 4,
    fontSize: 12,
    color: '#666',
  },
  uploadOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  uploadProgress: {
    backgroundColor: '#fff',
    padding: 24,
    borderRadius: 12,
    alignItems: 'center',
    margin: 20,
  },
  uploadText: {
    marginTop: 16,
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  uploadSubtext: {
    marginTop: 4,
    fontSize: 14,
    color: '#666',
  },
});

export default VideoScreen;