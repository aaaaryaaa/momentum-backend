// import React, { useEffect, useState } from 'react';
// import { useNavigation } from '@react-navigation/native';
// import {
//   View,
//   Text,
//   TouchableOpacity,
//   StyleSheet,
//   SafeAreaView,
//   ActivityIndicator,
//   Alert,
// } from 'react-native';
// import AsyncStorage from '@react-native-async-storage/async-storage';

// export default function HomeScreen({ onLogout }) {
//   const navigation = useNavigation();
//   const [userInfo, setUserInfo] = useState(null);
//   const [loading, setLoading] = useState(true);

//   const fetchUserInfo = async () => {
//     try {
//       const token = await AsyncStorage.getItem('token');
//       if (!token) {
//         Alert.alert('Session expired');
//         onLogout();
//         return;
//       }

//       const meResponse = await fetch('http://10.0.2.2:8000/me', {
//         headers: { Authorization: `Bearer ${token}` },
//       });

//       const data = await meResponse.json();

//       if (meResponse.ok) {
//         setUserInfo(data);
//       } else {
//         Alert.alert('Error', data.message || 'Failed to fetch user');
//         onLogout();
//       }
//     } catch (e) {
//       console.log('Error fetching user info:', e);
//       onLogout();
//     } finally {
//       setLoading(false);
//     }
//   };

//   const handleLogout = async () => {
//     await AsyncStorage.removeItem('token');
//     onLogout();
//   };

//   useEffect(() => {
//     fetchUserInfo();
//   }, []);

//   if (loading) {
//     return (
//       <SafeAreaView style={styles.container}>
//         <ActivityIndicator size="large" color="#007AFF" />
//       </SafeAreaView>
//     );
//   }

//   return (
//     <SafeAreaView style={styles.container}>
//       <View style={styles.header}>
//         <Text style={styles.welcomeText}>Welcome!</Text>
//         <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
//           <Text style={styles.logoutButtonText}>Logout</Text>
//         </TouchableOpacity>
//       </View>

//       <View style={styles.content}>
//         <Text style={styles.hiText}>
//           Hi {userInfo?.name || 'User'}! 👋
//         </Text>
//         <Text style={styles.emailText}>
//           {userInfo?.email || 'No email available'}
//         </Text>

//         <View style={styles.infoCard}>
//           <Text style={styles.cardTitle}>Your Profile</Text>
//           <Text style={styles.cardText}>Name: {userInfo?.name}</Text>
//           <Text style={styles.cardText}>Email: {userInfo?.email}</Text>
//         </View>

//         <View style={styles.comingSoonCard}>
//           <Text style={styles.cardTitle}>Coming Soon</Text>
//           <Text style={styles.cardText}>🎥 Video Upload</Text>
//           <Text style={styles.cardText}>💬 Chat Features</Text>
//           <Text style={styles.cardText}>📸 Instagram-like Features</Text>
//         </View>
//       </View>

//       <TouchableOpacity 
//         style={styles.chatButton} 
//         onPress={() => navigation.navigate('Chat')}
//       >
//         <Text style={styles.chatButtonText}>💬 Open Chat</Text>
//       </TouchableOpacity>
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
//     padding: 20,
//     backgroundColor: 'white',
//     borderBottomWidth: 1,
//     borderBottomColor: '#e0e0e0',
//   },
//   welcomeText: {
//     fontSize: 24,
//     fontWeight: 'bold',
//     color: '#333',
//   },
//   logoutButton: {
//     backgroundColor: '#FF3B30',
//     paddingHorizontal: 15,
//     paddingVertical: 8,
//     borderRadius: 8,
//   },
//   logoutButtonText: {
//     color: 'white',
//     fontWeight: 'bold',
//   },
//   content: {
//     flex: 1,
//     padding: 20,
//   },
//   hiText: {
//     fontSize: 32,
//     fontWeight: 'bold',
//     color: '#333',
//     marginBottom: 10,
//     textAlign: 'center',
//   },
//   emailText: {
//     fontSize: 16,
//     color: '#666',
//     marginBottom: 30,
//     textAlign: 'center',
//   },
//   infoCard: {
//     backgroundColor: 'white',
//     padding: 20,
//     borderRadius: 12,
//     marginBottom: 20,
//     shadowColor: '#000',
//     shadowOffset: {
//       width: 0,
//       height: 2,
//     },
//     shadowOpacity: 0.1,
//     shadowRadius: 4,
//     elevation: 3,
//   },
//   comingSoonCard: {
//     backgroundColor: '#E3F2FD',
//     padding: 20,
//     borderRadius: 12,
//     shadowColor: '#000',
//     shadowOffset: {
//       width: 0,
//       height: 2,
//     },
//     shadowOpacity: 0.1,
//     shadowRadius: 4,
//     elevation: 3,
//   },
//   cardTitle: {
//     fontSize: 20,
//     fontWeight: 'bold',
//     marginBottom: 10,
//     color: '#333',
//   },
//   cardText: {
//     fontSize: 16,
//     color: '#666',
//     marginBottom: 5,
//   },
//   chatButton: {
//     backgroundColor: '#007AFF',
//     padding: 15,
//     borderRadius: 12,
//     marginTop: 20,
//     alignItems: 'center',
//   },
//   chatButtonText: {
//     color: 'white',
//     fontSize: 18,
//     fontWeight: 'bold',
//   },
// });

import React, { useEffect, useState } from 'react';
import { useNavigation } from '@react-navigation/native';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  SafeAreaView,
  ActivityIndicator,
  Alert,
  ScrollView,
  Dimensions,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';

const { width } = Dimensions.get('window');

export default function HomeScreen({ onLogout }) {
  const navigation = useNavigation();
  const [userInfo, setUserInfo] = useState(null);
  const [loading, setLoading] = useState(true);

  const fetchUserInfo = async () => {
    try {
      const token = await AsyncStorage.getItem('token');
      if (!token) {
        Alert.alert('Session expired');
        onLogout();
        return;
      }

      const meResponse = await fetch('http://aarya.live:8000/me', {
        headers: { Authorization: `Bearer ${token}` },
      });

      const data = await meResponse.json();

      if (meResponse.ok) {
        setUserInfo(data);
      } else {
        Alert.alert('Error', data.message || 'Failed to fetch user');
        onLogout();
      }
    } catch (e) {
      console.log('Error fetching user info:', e);
      onLogout();
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    await AsyncStorage.removeItem('token');
    onLogout();
  };

  useEffect(() => {
    fetchUserInfo();
  }, []);

  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#007AFF" />
          <Text style={styles.loadingText}>Loading...</Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerContent}>
            <Text style={styles.welcomeText}>Welcome Back!</Text>
            <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
              <Text style={styles.logoutButtonText}>Logout</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* User Info Section */}
        <View style={styles.userSection}>
          <View style={styles.avatarContainer}>
            <View style={styles.avatar}>
              <Text style={styles.avatarText}>
                {userInfo?.name?.charAt(0)?.toUpperCase() || 'U'}
              </Text>
            </View>
          </View>
          <Text style={styles.hiText}>
            Hi {userInfo?.name || 'User'}! 👋
          </Text>
          <Text style={styles.emailText}>
            {userInfo?.email || 'No email available'}
          </Text>
        </View>

        {/* Quick Actions */}
        <View style={styles.quickActionsContainer}>
          <Text style={styles.sectionTitle}>Quick Actions</Text>
          <View style={styles.actionButtonsContainer}>
            <TouchableOpacity 
              style={[styles.actionButton, styles.chatButton]} 
              onPress={() => navigation.navigate('Chat')}
            >
              <View style={styles.actionButtonIcon}>
                <Text style={styles.actionButtonEmoji}>💬</Text>
              </View>
              <Text style={styles.actionButtonTitle}>Chat</Text>
              <Text style={styles.actionButtonSubtitle}>Start a conversation</Text>
            </TouchableOpacity>

            <TouchableOpacity 
              style={[styles.actionButton, styles.videoButton]} 
              onPress={() => navigation.navigate('Video')}
            >
              <View style={styles.actionButtonIcon}>
                <Text style={styles.actionButtonEmoji}>🎥</Text>
              </View>
              <Text style={styles.actionButtonTitle}>Videos</Text>
              <Text style={styles.actionButtonSubtitle}>Upload & manage videos</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Profile Card */}
        <View style={styles.profileCard}>
          <Text style={styles.cardTitle}>Your Profile</Text>
          <View style={styles.profileInfo}>
            <View style={styles.profileRow}>
              <Text style={styles.profileLabel}>Name:</Text>
              <Text style={styles.profileValue}>{userInfo?.name || 'Not set'}</Text>
            </View>
            <View style={styles.profileRow}>
              <Text style={styles.profileLabel}>Email:</Text>
              <Text style={styles.profileValue}>{userInfo?.email || 'Not set'}</Text>
            </View>
          </View>
        </View>

        {/* Coming Soon Card */}
        <View style={styles.comingSoonCard}>
          <Text style={styles.cardTitle}>Coming Soon</Text>
          <View style={styles.featureList}>
            <View style={styles.featureItem}>
              <Text style={styles.featureEmoji}>📸</Text>
              <Text style={styles.featureText}>Instagram-like Stories</Text>
            </View>
            <View style={styles.featureItem}>
              <Text style={styles.featureEmoji}>🔔</Text>
              <Text style={styles.featureText}>Push Notifications</Text>
            </View>
            <View style={styles.featureItem}>
              <Text style={styles.featureEmoji}>👥</Text>
              <Text style={styles.featureText}>Social Features</Text>
            </View>
            <View style={styles.featureItem}>
              <Text style={styles.featureEmoji}>⚡</Text>
              <Text style={styles.featureText}>Real-time Updates</Text>
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  header: {
    backgroundColor: '#ffffff',
    paddingTop: 10,
    paddingBottom: 20,
    borderBottomLeftRadius: 25,
    borderBottomRightRadius: 25,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 5,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
  },
  welcomeText: {
    fontSize: 26,
    fontWeight: 'bold',
    color: '#2c3e50',
  },
  logoutButton: {
    backgroundColor: '#e74c3c',
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
    shadowColor: '#e74c3c',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 4,
    elevation: 3,
  },
  logoutButtonText: {
    color: 'white',
    fontWeight: '600',
    fontSize: 14,
  },
  userSection: {
    alignItems: 'center',
    paddingVertical: 30,
    paddingHorizontal: 20,
  },
  avatarContainer: {
    marginBottom: 15,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#3498db',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#3498db',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 5,
  },
  avatarText: {
    fontSize: 32,
    fontWeight: 'bold',
    color: 'white',
  },
  hiText: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 5,
  },
  emailText: {
    fontSize: 16,
    color: '#7f8c8d',
  },
  quickActionsContainer: {
    paddingHorizontal: 20,
    marginBottom: 25,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 15,
  },
  actionButtonsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 15,
  },
  actionButton: {
    flex: 1,
    padding: 20,
    borderRadius: 20,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 5,
  },
  chatButton: {
    backgroundColor: '#3498db',
  },
  videoButton: {
    backgroundColor: '#e74c3c',
  },
  actionButtonIcon: {
    marginBottom: 10,
  },
  actionButtonEmoji: {
    fontSize: 32,
  },
  actionButtonTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: 'white',
    marginBottom: 5,
  },
  actionButtonSubtitle: {
    fontSize: 12,
    color: 'rgba(255, 255, 255, 0.8)',
    textAlign: 'center',
  },
  profileCard: {
    backgroundColor: 'white',
    marginHorizontal: 20,
    marginBottom: 20,
    padding: 20,
    borderRadius: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 3,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 15,
  },
  profileInfo: {
    gap: 10,
  },
  profileRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  profileLabel: {
    fontSize: 14,
    color: '#7f8c8d',
    fontWeight: '500',
  },
  profileValue: {
    fontSize: 14,
    color: '#2c3e50',
    fontWeight: '600',
    flex: 1,
    textAlign: 'right',
  },
  comingSoonCard: {
    backgroundColor: '#ecf0f1',
    marginHorizontal: 20,
    marginBottom: 30,
    padding: 20,
    borderRadius: 20,
    borderWidth: 2,
    borderColor: '#bdc3c7',
    borderStyle: 'dashed',
  },
  featureList: {
    gap: 12,
  },
  featureItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  featureEmoji: {
    fontSize: 20,
  },
  featureText: {
    fontSize: 14,
    color: '#7f8c8d',
    flex: 1,
  },
});