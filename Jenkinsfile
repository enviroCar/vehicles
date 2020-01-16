node {
  def database
  def service

  stage('Clone repository') {
    checkout scm
  }

  stage('Build database') {
    database = docker.build("envirocar/vehicles-db", "db/")
  }
  
  stage('Build service') {
    service = docker.build("envirocar/vehicles")
  }

  stage('Push images') {
    docker.withRegistry('http://registry:5000') {
      database.push("latest")
      service.push("latest")
    }
  }
}