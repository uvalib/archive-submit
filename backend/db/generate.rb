ainsert = "insert into accessions(id,identifier,user_id,description,accession_type,created_at) values "
ginsert = "insert into accession_genres (accession_id,genre_id) values "
dinsert = "insert into digital_accessions (accession_id,upload_size,description) values "
start_id=10
identifier="bj4u8ucrtr351n8a4q" ## missing 2 digits

accessions = []
genres = []
das = []
while start_id<100 

   accessions << "(#{start_id},'#{identifier}%02d',1,'This is fake description number %d','new','%s')" % [start_id, start_id,Time.now.strftime("%F %T")]
   genre_id =  Random.rand(11)+1
   genres << "(#{start_id},#{genre_id})"
   das << "(#{start_id},0,'Fake digital accession number #{start_id}')" 
   start_id += 1
end

puts ainsert
puts "#{accessions.join(",\n")};"
puts
puts ginsert
puts "#{genres.join(", ")};"
puts
puts dinsert
puts "#{das.join(",\n")};"